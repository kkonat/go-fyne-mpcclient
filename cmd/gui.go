package main

import (
	"fmt"
	"time"

	fe "remotecc/fynextensions"

	f2 "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ControlCenterPanelGUI struct {
	statusL              *canvas.Text
	artist, album, track *canvas.Text
	bPower               *widget.Button
	vol                  binding.Float
	elapsed              binding.Float
	lastVol              TrackVolume
	volLastChngd         time.Time
	prgrs                *fe.TappableProgressBar
	storedStatusText     string
	IPaddrs              binding.String
}

func newControlCenterAppGUI(w f2.Window, app *ControlCenterApp) *ControlCenterPanelGUI {

	c := &ControlCenterPanelGUI{

		artist:  fe.NewText("Artist:", 12),
		album:   fe.NewText("Album:", 12),
		track:   fe.NewText("Trk:", 12),
		statusL: fe.NewText("Ready", 10),
		vol:     binding.NewFloat(),
		elapsed: binding.NewFloat(),
		IPaddrs: binding.NewString()}

	fe.NewText("Volume:", 10)
	c.IPaddrs.Set("192.168.0.95:6600")

	bInputPlayer := widget.NewButton("Player", func() {
		c.setStatusText("switched to Player", 3)
		app.hw.Request("ctrl", "deq_input_coaxial")
	})

	bInputTV := widget.NewButton("TV", func() {
		c.setStatusText("switched to TV", 3)
		app.hw.Request("ctrl", "deq_input_optical")
	})

	c.bPower = widget.NewButton("Power", func() {
		c.togglePower(app)
	})

	bShtDn := widget.NewButton("Shutdown", func() {
		c.setStatusText("Shutting down...", 0)
		app.hw.Request("ctrl", "server_poweroff")
	})

	slider := widget.NewSliderWithData(0, 100, c.vol)
	slider.Orientation = widget.Orientation(f2.OrientationVerticalUpsideDown)
	slider.Move(f2.NewPos(0, 20))

	c.lastVol = app.state.volume
	c.volLastChngd = time.Now()

	// volume slider
	// slider dragging
	slider.OnChanged = func(val float64) {
		if time.Since(c.volLastChngd).Milliseconds() > 100 {
			c.volLastChngd = time.Now()
			v := int(val)
			app.hw.Request("mpd", fmt.Sprintf("setvol %d", v))
		}
	}

	c.prgrs = fe.NewTappableProgressBarWithData(c.elapsed)
	c.prgrs.OnTapped = func(v float64) {
		seek := float64(app.state.track.duration) * v
		song := app.state.track.song
		app.hw.Request("mpd", fmt.Sprintf("seek %d %d", song, int(seek)))
	}
	c.prgrs.Max = float64(app.state.track.duration)
	c.prgrs.TextFormatter = func() string {
		return trkTimeToString(float32(c.prgrs.Value))
	}

	// menu := f2.NewMainMenu(f2.NewMenu("...",
	// 	f2.NewMenuItem("Config", func() {}),
	// 	f2.NewMenuItem("Show hw controls", func() {})))

	conPlayer := container.NewBorder(
		nil,
		container.NewVBox(widget.NewSeparator(), c.statusL),
		nil,
		slider,
		container.NewVBox(
			c.artist, c.album, c.track, c.prgrs,
			container.NewGridWithColumns(5,
				widget.NewButtonWithIcon("", theme.MediaSkipPreviousIcon(), func() {
					app.hw.Request("mpd", "previous")
					c.setStatusText("skip back", 1)
				}),
				widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {
					app.hw.Request("mpd", "play")
				}),
				widget.NewButtonWithIcon("", theme.MediaPauseIcon(), func() {
					app.hw.Request("mpd", "pause")
				}),
				widget.NewButtonWithIcon("", theme.MediaStopIcon(), func() {
					app.hw.Request("mpd", "stop")
				}),
				widget.NewButtonWithIcon("", theme.MediaSkipNextIcon(), func() {
					app.hw.Request("mpd", "next")
					c.setStatusText("skip next", 1)
				})),
		),
	)
	conHW := container.NewBorder(nil, container.NewVBox(widget.NewSeparator(), c.statusL), nil, nil,
		container.NewVBox(container.NewGridWithColumns(2, bInputPlayer, bInputTV),
			c.bPower,
			bShtDn),
	)
	conSettings := container.NewBorder(nil, container.NewVBox(widget.NewSeparator(), c.statusL), nil, nil,
		container.NewVBox(
			widget.NewLabel("Settings"),
			widget.NewSeparator(),
			widget.NewLabel("IP:"),
			widget.NewEntryWithData(c.IPaddrs)),
	)
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("", theme.MediaMusicIcon(), conPlayer),
		container.NewTabItemWithIcon("", theme.ComputerIcon(), conHW),
		container.NewTabItemWithIcon("", theme.SettingsIcon(), conSettings),
	)

	tabs.SetTabLocation(container.TabLocationTop)

	w.SetContent(tabs)
	// w.SetMainMenu(menu)
	return c
}
func (c ControlCenterPanelGUI) updatePowerBbutton(powerOn bool) {
	if powerOn {
		c.bPower.SetText("Power off")
	} else {
		c.bPower.SetText("Power on")
	}
}
func (c ControlCenterPanelGUI) updateOnlineStatus(online bool, ps PlayStatus) {
	if !online {
		c.setStatusText("Offline. Waiting for connection...", 0)
	} else {
		c.UpdatePlayStatus(ps)
	}
}

func (c ControlCenterPanelGUI) togglePower(app *ControlCenterApp) {
	var state bool
	var err error
	if err = app.hw.togglePower(); err == nil {
		if state, err = app.state.getHWState(); err == nil {
			if state {
				c.setStatusText("powered off", 3)
			} else {
				c.setStatusText("powered on", 3)
			}
			c.updatePowerBbutton(state)
		}
	}
}
func (c ControlCenterPanelGUI) setStatusText(newText string, howLong int) {
	if howLong == 0 {
		c.statusL.Text = newText
		c.statusL.Refresh()
		return
	}

	c.storedStatusText = c.statusL.Text
	c.statusL.Text = newText
	time.AfterFunc(time.Second*time.Duration(howLong),
		func() {
			c.statusL.Text = c.storedStatusText
			c.statusL.Refresh()
		})
	c.statusL.Refresh()
}
func (c *ControlCenterPanelGUI) UpdatePlayStatus(s PlayStatus) {
	switch s {
	case playing:
		c.setStatusText("Playing...", 0)
	case stopped:
		c.setStatusText("Stopped", 0)
	case paused:
		c.setStatusText("Paused", 0)
	}
}
func (c *ControlCenterPanelGUI) updateTrackDetails(ti *TrackInfo) {
	c.album.Text = "Album: " + ti.album
	c.artist.Text = "Artist: " + ti.artist
	c.track.Text = ti.track
	c.album.Refresh()
	c.artist.Refresh()
	c.track.Refresh()
	c.prgrs.Max = float64(ti.duration)
}

func (c *ControlCenterPanelGUI) updateTrackElapsedTime(elTime TrackTime) {
	c.elapsed.Set(float64(elTime))
}

func (c *ControlCenterPanelGUI) updateVolume(v TrackVolume) {
	c.vol.Set(float64(v))
	c.lastVol = v
}

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
	prgrs                *widget.ProgressBar
}

func newControlCenterAppGUI(w f2.Window, app *ControlCenterApp) *ControlCenterPanelGUI {

	c := &ControlCenterPanelGUI{

		artist:  fe.NewText("Artist:", 12),
		album:   fe.NewText("Album:", 12),
		track:   fe.NewText("Trk:", 12),
		statusL: fe.NewText("Ready", 10),
		vol:     binding.NewFloat(),
		elapsed: binding.NewFloat()}

	fe.NewText("Volume:", 10)

	bInputPlayer := widget.NewButton("Player", func() {
		c.setStatusText("switched to Player")
		app.hw.Request("ctrl", "deq_input_coaxial")
	})

	bInputTV := widget.NewButton("TV", func() {
		c.setStatusText("switched to TV")
		app.hw.Request("ctrl", "deq_input_optical")
	})

	c.bPower = widget.NewButton("Power", func() {
		c.togglePower(app)
	})

	bShtDn := widget.NewButton("Shutdown", func() {
		c.setStatusText("Shutting down...")
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

	c.prgrs = widget.NewProgressBarWithData(c.elapsed)
	c.prgrs.Max = float64(app.state.track.duration)
	c.prgrs.TextFormatter = func() string {
		return trkTimeToString(float32(c.prgrs.Value))
	}

	con := container.NewBorder(
		nil,
		container.NewVBox(widget.NewSeparator(), c.statusL),
		nil,
		slider,
		container.NewVBox(
			c.artist, c.album, c.track, c.prgrs, widget.NewSeparator(),
			container.NewGridWithColumns(2, bInputPlayer, bInputTV),
			container.NewGridWithColumns(5,
				widget.NewButtonWithIcon("", theme.MediaSkipPreviousIcon(), func() {
					app.hw.Request("mpd", "previous")
					c.setStatusText("skip back")
				}),
				widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {
					app.hw.Request("mpd", "play")
					c.setStatusText("playing")
				}),
				widget.NewButtonWithIcon("", theme.MediaPauseIcon(), func() {
					app.hw.Request("mpd", "pause")
					c.setStatusText("paused")
				}),
				widget.NewButtonWithIcon("", theme.MediaStopIcon(), func() {
					app.hw.Request("mpd", "stop")
					c.setStatusText("stopped")
				}),
				widget.NewButtonWithIcon("", theme.MediaSkipNextIcon(), func() {
					app.hw.Request("mpd", "next")
					c.setStatusText("skip next")
				})),
			c.bPower,
			bShtDn),
	)

	w.SetContent(con)

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
		c.setStatusText("Offline. Waiting for connection...")
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
				c.setStatusText("powered on")
			} else {
				c.setStatusText("powered off")
			}
			c.updatePowerBbutton(state)
		}
	}
}
func (c ControlCenterPanelGUI) setStatusText(t string) {
	c.statusL.Text = t
	c.statusL.Refresh()
}
func (c *ControlCenterPanelGUI) UpdatePlayStatus(s PlayStatus) {
	switch s {
	case playing:
		c.setStatusText("Playing...")
	case stopped:
		c.setStatusText("Stopped")
	case paused:
		c.setStatusText("Paused")
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

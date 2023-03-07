package gui

import (
	"fmt"
	"time"

	fe "remotecc/cmd/fynextensions"
	hw "remotecc/cmd/hwinterface"
	"remotecc/cmd/state"

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
	lastVol              state.TrackVolume
	volLastChngd         time.Time
	prgrs                *fe.TappableProgressBar
	storedStatusText     string
	statusPopupStart     time.Time

	stateStream chan any

	Hw    *hw.HWInterface
	State *state.PlayerState
}

func New(w f2.Window, stream chan any, State *state.PlayerState, Hw *hw.HWInterface) *ControlCenterPanelGUI {

	c := &ControlCenterPanelGUI{
		stateStream: stream,
		Hw:          Hw, State: State,
		artist:  fe.NewText("Artist:", 12),
		album:   fe.NewText("Album:", 12),
		track:   fe.NewText("Trk:", 12),
		statusL: fe.NewText("Ready", 10),
		vol:     binding.NewFloat(),
		elapsed: binding.NewFloat()}

	fe.NewText("Volume:", 10)

	bInputPlayer := widget.NewButton("Player", func() {
		c.SetStatusText("switched to Player", 3)
		Hw.Request("ctrl", "deq_input_coaxial")
	})

	bInputTV := widget.NewButton("TV", func() {
		c.SetStatusText("switched to TV", 3)
		Hw.Request("ctrl", "deq_input_optical")
	})

	c.bPower = widget.NewButton("Power", func() {
		c.TogglePower()
	})

	bShtDn := widget.NewButton("Shutdown", func() {
		c.SetStatusText("Shutting down...", 0)
		Hw.Request("ctrl", "server_poweroff")
	})

	slider := widget.NewSliderWithData(0, 100, c.vol)
	slider.Orientation = widget.Orientation(f2.OrientationVerticalUpsideDown)
	slider.Move(f2.NewPos(0, 20))

	c.lastVol = State.Volume
	c.volLastChngd = time.Now()

	// volume slider
	// slider dragging
	slider.OnChanged = func(val float64) {
		if time.Since(c.volLastChngd).Milliseconds() > 100 {
			c.volLastChngd = time.Now()
			v := int(val)
			Hw.Request("mpd", fmt.Sprintf("setvol %d", v))
		}
	}

	c.prgrs = fe.NewTappableProgressBarWithData(c.elapsed)
	c.prgrs.OnTapped = func(v float64) {
		seek := float64(State.Track.Duration) * v
		song := State.Track.Song
		Hw.Request("mpd", fmt.Sprintf("seek %d %d", song, int(seek)))
	}
	c.prgrs.Max = float64(State.Track.Duration)
	c.prgrs.TextFormatter = func() string {
		return state.TrkTimeToString(float32(c.prgrs.Value))
	}
	bConf := widget.NewButton("Config", func() {})

	con := container.NewBorder(
		nil,
		container.NewVBox(widget.NewSeparator(), c.statusL),
		nil,
		slider,
		container.NewVBox(
			bConf,
			c.artist, c.album, c.track, c.prgrs,
			container.NewGridWithColumns(5,
				widget.NewButtonWithIcon("", theme.MediaSkipPreviousIcon(), func() {
					Hw.Request("mpd", "previous")
					c.SetStatusText("skip back", 1)
				}),
				widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {
					Hw.Request("mpd", "play")
				}),
				widget.NewButtonWithIcon("", theme.MediaPauseIcon(), func() {
					Hw.Request("mpd", "pause")
				}),
				widget.NewButtonWithIcon("", theme.MediaStopIcon(), func() {
					Hw.Request("mpd", "stop")
				}),
				widget.NewButtonWithIcon("", theme.MediaSkipNextIcon(), func() {
					Hw.Request("mpd", "next")
					c.SetStatusText("skip next", 1)
				})),
			widget.NewSeparator(),
			container.NewGridWithColumns(2, bInputPlayer, bInputTV),
			c.bPower,
			bShtDn,
		),
	)

	w.SetContent(con)

	return c
}
func (c ControlCenterPanelGUI) UpdatePowerBbutton(powerOn bool) {
	if powerOn {
		c.bPower.SetText("Power off")
	} else {
		c.bPower.SetText("Power on")
	}
}
func (c ControlCenterPanelGUI) UpdateOnlineStatus(online bool, ps state.PlayStatus) {
	if !online {
		c.SetStatusText("Offline. Waiting for connection...", 0)
	} else {
		c.UpdatePlayStatus(ps)
	}
}

func (c ControlCenterPanelGUI) TogglePower() {
	var state bool
	var err error
	if err = c.Hw.TogglePower(); err == nil {
		if state, err = c.State.GetHWState(); err == nil {
			if state {
				c.SetStatusText("powered off", 3)
			} else {
				c.SetStatusText("powered on", 3)
			}
			c.UpdatePowerBbutton(state)
		}
	}
}
func (c ControlCenterPanelGUI) SetStatusText(newText string, howLong int) {
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
func (c *ControlCenterPanelGUI) UpdatePlayStatus(s state.PlayStatus) {
	switch s {
	case state.Playing:
		c.SetStatusText("Playing...", 0)
	case state.Stopped:
		c.SetStatusText("Stopped", 0)
	case state.Paused:
		c.SetStatusText("Paused", 0)
	}
}
func (c *ControlCenterPanelGUI) UpdateTrackDetails(ti *state.TrackInfo) {
	c.album.Text = "Album: " + ti.Album
	c.artist.Text = "Artist: " + ti.Artist
	c.track.Text = ti.Track
	c.album.Refresh()
	c.artist.Refresh()
	c.track.Refresh()
	c.prgrs.Max = float64(ti.Duration)
}

func (c *ControlCenterPanelGUI) UpdateTrackElapsedTime(elTime state.TrackTime) {
	c.elapsed.Set(float64(elTime))
}

func (c *ControlCenterPanelGUI) UpdateVolume(v state.TrackVolume) {
	c.vol.Set(float64(v))
	c.lastVol = v
}

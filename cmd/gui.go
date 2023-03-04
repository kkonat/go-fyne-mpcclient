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
	vol                  binding.Float
	elapsed              binding.Float
	lastVol              TrackVolume
	volLastChngd         time.Time
	prgrs                *widget.ProgressBar
}

func newControlCenterPanelGUI(w f2.Window, mpc *ControlCenterApp) *ControlCenterPanelGUI {

	c := new(ControlCenterPanelGUI)

	c.artist = fe.NewText("Artist:", 12)
	c.album = fe.NewText("Album:", 12)
	c.track = fe.NewText("Trk:", 12)

	c.updateTrackDetails(&mpc.state.track)

	c.statusL = fe.NewText("Ready", 10)

	c.vol = binding.NewFloat()
	c.elapsed = binding.NewFloat()

	fe.NewText("Volume:", 10)

	bInputPlayer := widget.NewButton("Player", func() {
		c.StatusText("switched to Player")
		mpc.ctrlClient.ask("deq_input_coaxial")
	})

	bInputTV := widget.NewButton("TV", func() {
		c.StatusText("switched to TV")
		mpc.ctrlClient.ask("deq_input_optical")
	})

	bPower := widget.NewButton("Power", c.tapPower)
	bShtDn := widget.NewButton("Shutdown", c.tapShtdn)

	sl := widget.NewSliderWithData(0, 100, c.vol)
	sl.Orientation = widget.Orientation(f2.OrientationVerticalUpsideDown)
	sl.Move(f2.NewPos(0, 20))

	c.lastVol = mpc.state.volume
	c.volLastChngd = time.Now()

	// volume slider
	// slider dragging
	sl.OnChanged = func(val float64) {
		elapsed := time.Since(c.volLastChngd).Milliseconds()

		if elapsed > 100 {
			c.volLastChngd = time.Now()
			v := int(val)
			mpc.mpdClient.ask(fmt.Sprintf("setvol %d", v))
		}
	}

	c.vol.Set(float64(mpc.state.volume))

	c.prgrs = widget.NewProgressBarWithData(c.elapsed)
	c.prgrs.Max = float64(mpc.state.track.duration)
	c.prgrs.TextFormatter = func() string {
		return trkTimeToString(float32(c.prgrs.Value))
	}

	con := container.NewBorder(
		nil,
		container.NewVBox(widget.NewSeparator(), c.statusL),
		nil,
		sl,
		container.NewVBox(
			c.artist, c.album, c.track, c.prgrs, widget.NewSeparator(),
			container.NewGridWithColumns(2, bInputPlayer, bInputTV),
			container.NewGridWithColumns(5,
				widget.NewButtonWithIcon("", theme.MediaSkipPreviousIcon(), func() { mpc.mpdClient.ask("previous") }),
				widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() { mpc.mpdClient.ask("play") }),
				widget.NewButtonWithIcon("", theme.MediaPauseIcon(), func() { mpc.mpdClient.ask("pause") }),
				widget.NewButtonWithIcon("", theme.MediaStopIcon(), func() { mpc.mpdClient.ask("stop") }),
				widget.NewButtonWithIcon("", theme.MediaSkipNextIcon(), func() { mpc.mpdClient.ask("next") })),
			bPower, bShtDn),
	)

	w.SetContent(con)

	return c
}

func (c ControlCenterPanelGUI) StatusText(t string) {
	c.statusL.Text = t
	c.statusL.Refresh()
}
func (c ControlCenterPanelGUI) chkPowerStatus() bool {
	// mpc.ctrlSrv.sendCtrlCmd( "check_extpower")
	// return strings.Split(s[0], ": ")[1] == "1"
	return true
}

func (c ControlCenterPanelGUI) tapPower() {
	if c.chkPowerStatus() {
		c.StatusText("powered off")
	} else {
		c.StatusText("powered on")
	}
	// sendCtrlCmd(ctrlSrv, "extpower_toggle")
}
func (c ControlCenterPanelGUI) tapShtdn() {
	c.StatusText("Shutting down...")
	// sendCtrlCmd(ctrlSrv, "server_poweroff")
}

func (c *ControlCenterPanelGUI) updateTrackDetails(ti *TrackInfo) {
	c.album.Text = "Album: " + ti.album
	c.album.Refresh()
	c.artist.Text = "Artist: " + ti.artist
	c.artist.Refresh()
	c.track.Text = ti.track
	c.track.Refresh()
	if c.prgrs != nil {
		c.prgrs.Max = float64(ti.duration)
	}
}

func (c *ControlCenterPanelGUI) updateTrackElapsedTime(elTime TrackTime) {
	c.elapsed.Set(float64(elTime))
}

func (c *ControlCenterPanelGUI) updateVolume(v TrackVolume) {
	if v != c.lastVol && time.Since(c.volLastChngd).Seconds() > 2 {
		c.vol.Set(float64(v))
		c.lastVol = v
	}
}
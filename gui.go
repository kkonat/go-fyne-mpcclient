package main

import (
	"fmt"
	"time"

	f2 "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ControlCenterPanel struct {
	statusL                    *canvas.Text
	artist, album, track, time *canvas.Text
	vol                        binding.Float
}

func newControlCenterPanelGUI(w f2.Window, mpc *MpcClient) *ControlCenterPanel {

	c := new(ControlCenterPanel)
	fg := theme.ForegroundColor()

	c.artist = canvas.NewText("Artist:", fg)
	c.artist.TextSize = 12
	c.album = canvas.NewText("Album:", fg)
	c.album.TextSize = 12
	c.track = canvas.NewText("tr:", fg)
	c.track.TextSize = 12
	c.time = canvas.NewText("Time:", fg)
	c.time.TextSize = 12

	c.updateTrackDetails(mpc.ps.track)

	// c.statusL = widget.NewLabel("ready")
	c.statusL = canvas.NewText("Ready", fg)
	c.statusL.TextSize = 10
	c.vol = binding.NewFloat()
	volLvl := canvas.NewText("Volume:", fg)
	volLvl.TextSize = 10
	bInputPlayer := widget.NewButton("Player", func() {
		c.StatusText("switched to Player")
		// sendCtrlCmd(ctrlSrv, "deq_input_coaxial")
	})
	bInputTV := widget.NewButton("TV", func() {
		c.StatusText("switched to TV")
		// sendCtrlCmd(ctrlSrv, "deq_input_optical")
	})
	bPower := widget.NewButton("Power", c.tapPower)
	bShtDn := widget.NewButton("Shutdown", c.tapShtdn)

	sl := widget.NewSliderWithData(0, 100, c.vol)
	sl.Orientation = widget.Orientation(f2.OrientationVerticalUpsideDown)
	sl.Move(f2.NewPos(0, 20))

	// lastVol := ps.volume
	volLastChngd := time.Now()

	// volume slider
	// slider dragging
	sl.OnChanged = func(val float64) {
		elapsed := time.Since(volLastChngd).Milliseconds()

		if elapsed > 100 {
			volLastChngd = time.Now()
			v := int(val)
			mpc.sendCtrlCmd("setvol " + fmt.Sprintf("%d", v))
		}
	}

	c.vol.Set(float64(mpc.ps.volume))

	prgrs := widget.NewProgressBar()
	trlen := mpc.ps.track.duration
	prgrs.Min = 0
	prgrs.Max = float64(trlen)

	con := container.NewBorder(
		nil,
		container.NewVBox(widget.NewSeparator(), c.statusL),
		nil,
		sl,
		container.NewVBox(
			c.artist, c.album, c.track, c.time, prgrs, widget.NewSeparator(),
			container.NewGridWithColumns(2, bInputPlayer, bInputTV),
			container.NewGridWithColumns(5,
				widget.NewButtonWithIcon("", theme.MediaSkipPreviousIcon(), func() { mpc.sendCtrlCmd("previous") }),
				widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() { mpc.sendCtrlCmd("play") }),
				widget.NewButtonWithIcon("", theme.MediaPauseIcon(), func() { mpc.sendCtrlCmd("pause") }),
				widget.NewButtonWithIcon("", theme.MediaStopIcon(), func() { mpc.sendCtrlCmd("stop") }),
				widget.NewButtonWithIcon("", theme.MediaSkipNextIcon(), func() { mpc.sendCtrlCmd("next") })),
			bPower, bShtDn),
	)

	w.SetContent(con)

	return c
}

func (c ControlCenterPanel) StatusText(t string) {
	c.statusL.Text = t
	c.statusL.Refresh()
}
func (c ControlCenterPanel) chkPowerStatus() bool {
	// s := sendCtrlCmd(ctrlSrv, "check_extpower")
	// return strings.Split(s[0], ": ")[1] == "1"
	return true
}

func (c ControlCenterPanel) tapPower() {
	if c.chkPowerStatus() {
		c.StatusText("powered off")
	} else {
		c.StatusText("powered on")
	}
	// sendCtrlCmd(ctrlSrv, "extpower_toggle")
}
func (c ControlCenterPanel) tapShtdn() {
	c.StatusText("Shutting down...")
	// sendCtrlCmd(ctrlSrv, "server_poweroff")
}

func (c *ControlCenterPanel) updateTrackDetails(ti TrackInfo) {
	c.album.Text = "Album: " + ti.album
	c.album.Refresh()
	c.artist.Text = "Artist: " + ti.artist
	c.artist.Refresh()
	c.track.Text = ti.track
	c.track.Refresh()
}

func (c *ControlCenterPanel) updateVolume(vol TrackVolume) {
	// if err == nil && v != lastVol && time.Since(volLastChngd).Seconds() > 2 {
	// 	c.vol.Set(float64(v))
	// 	lastVol = v
	// }
}

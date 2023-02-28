package main

import (
	"fmt"
	"strings"
	"time"

	f2 "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const ctrlSrv = "192.168.0.95:1025"
const mpdSrv = "192.168.0.95:6600"

type cc struct {
	statusL                    *canvas.Text
	artist, album, track, time *canvas.Text
	vol                        binding.Float
}

func newCc(w f2.Window) *cc {
	c := new(cc)
	fg := theme.ForegroundColor()

	c.artist = canvas.NewText("Artist:", fg)
	c.artist.TextSize = 12
	c.album = canvas.NewText("Album:", fg)
	c.album.TextSize = 12
	c.track = canvas.NewText("tr:", fg)
	c.track.TextSize = 12
	c.time = canvas.NewText("Time:", fg)
	c.time.TextSize = 12

	c.updateTrackDetails()

	// c.statusL = widget.NewLabel("ready")
	c.statusL = canvas.NewText("Ready", fg)
	c.statusL.TextSize = 10
	c.vol = binding.NewFloat()
	volLvl := canvas.NewText("Volume:", fg)
	volLvl.TextSize = 10
	bInputPlayer := widget.NewButton("Player", func() {
		c.StatusText("switched to Player")
		sendCtrlCmd(ctrlSrv, "deq_input_coaxial")
	})
	bInputTV := widget.NewButton("TV", func() {
		c.StatusText("switched to TV")
		sendCtrlCmd(ctrlSrv, "deq_input_optical")
	})
	bPower := widget.NewButton("Power", c.tapPower)
	bShtDn := widget.NewButton("Shutdown", c.tapShtdn)

	sl := widget.NewSliderWithData(0, 100, c.vol)
	sl.Orientation = widget.Orientation(f2.OrientationVerticalUpsideDown)
	sl.Move(f2.NewPos(0, 20))

	lastVol, _ := getVolume()
	volLastChngd := time.Now()

	// volume slider

	// periodic update

	go func() {
		for {
			v, err := getVolume()
			// if new value and slider was dragged over 2 secs ago
			if err == nil && v != lastVol && time.Since(volLastChngd).Seconds() > 2 {
				c.vol.Set(float64(v))
				lastVol = v
			}
			// ..., sleep, repeat
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// slider dragging
	sl.OnChanged = func(val float64) {
		elapsed := time.Since(volLastChngd).Milliseconds()

		if elapsed > 100 {
			volLastChngd = time.Now()
			v := int(val)
			sendCtrlCmd(mpdSrv, "setvol "+fmt.Sprintf("%d", v))
		}
	}

	vol, err := getVolume()
	if err == nil {
		c.vol.Set(float64(vol))
	}
	prgrs := widget.NewProgressBar()

	trlen, err := getTrackLen()
	if err == nil {
		prgrs.Min = 0
		prgrs.Max = float64(trlen)
		fmt.Println("trlen=", trlen)
	}

	con := container.NewBorder(
		nil,
		container.NewVBox(widget.NewSeparator(), c.statusL),
		nil,
		sl,
		container.NewVBox(
			c.artist, c.album, c.track, c.time, prgrs, widget.NewSeparator(),
			container.NewGridWithColumns(2, bInputPlayer, bInputTV),
			container.NewGridWithColumns(5,
				widget.NewButtonWithIcon("", theme.MediaSkipPreviousIcon(), func() { sendCtrlCmd(mpdSrv, "previous") }),
				widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() { sendCtrlCmd(mpdSrv, "play") }),
				widget.NewButtonWithIcon("", theme.MediaPauseIcon(), func() { sendCtrlCmd(mpdSrv, "pause") }),
				widget.NewButtonWithIcon("", theme.MediaStopIcon(), func() { sendCtrlCmd(mpdSrv, "stop") }),
				widget.NewButtonWithIcon("", theme.MediaSkipNextIcon(), func() { sendCtrlCmd(mpdSrv, "next") })),
			bPower, bShtDn),
	)

	// con := container.NewVBox(container.NewHBox(container.NewVBox(
	// 	c.artist, c.album, c.track, widget.NewSeparator(),
	// 	container.NewGridWithColumns(2, b1, b2),
	// 	container.NewGridWithColumns(5,
	// 		widget.NewButtonWithIcon("", theme.MediaSkipPreviousIcon(), func() { sendCtrlCmd(mpdSrv, "prev") }),
	// 		widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() { sendCtrlCmd(mpdSrv, "play") }),
	// 		widget.NewButtonWithIcon("", theme.MediaPauseIcon(), func() { sendCtrlCmd(mpdSrv, "pause") }),
	// 		widget.NewButtonWithIcon("", theme.MediaStopIcon(), func() { sendCtrlCmd(mpdSrv, "stop") }),
	// 		widget.NewButtonWithIcon("", theme.MediaSkipNextIcon(), func() { sendCtrlCmd(mpdSrv, "next") })),
	// 	bPower,
	// ), sl),
	// 	widget.NewSeparator(),
	// 	c.statusL,
	// )

	w.SetContent(con)

	return c
}

func (c cc) StatusText(t string) {
	c.statusL.Text = t
	c.statusL.Refresh()
}
func (c cc) chkPowerStatus() bool {
	s := sendCtrlCmd(ctrlSrv, "check_extpower")
	return strings.Split(s[0], ": ")[1] == "1"
}

func (c cc) tapPower() {
	if c.chkPowerStatus() {
		c.StatusText("powered off")
	} else {
		c.StatusText("powered on")
	}
	sendCtrlCmd(ctrlSrv, "extpower_toggle")
}
func (c cc) tapShtdn() {
	c.StatusText("Shutting down...")
	sendCtrlCmd(ctrlSrv, "server_poweroff")
}

func (c *cc) updateTrackDetails() {
	resp := sendCtrlCmd(mpdSrv, "currentsong")
	c.album.Text = "Album: " + extract(resp, "Album:")
	c.album.Refresh()
	c.artist.Text = "Artist: " + extract(resp, "Artist")
	c.artist.Refresh()
	c.track.Text = extract(resp, "Track:") + " " + extract(resp, "Title")
	c.track.Refresh()

}

func main() {
	a := app.New()
	w := a.NewWindow("Remote Control Center")

	cc := newCc(w)

	oldHash := getTrackDataHash()
	go func() {
		for {
			hash := getTrackDataHash()
			if oldHash != hash {
				cc.updateTrackDetails()
				oldHash = hash
			}
			time.Sleep(1 * time.Second)
		}
	}()
	w.ShowAndRun()
}

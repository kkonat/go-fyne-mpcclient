package gui

import (
	"fmt"
	"remotecc/cmd/state"
	"time"

	f2 "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type VolumeSlider struct {
	volSlider    *widget.Slider
	vol          binding.Float
	LastVol      state.TrackVolume
	volPremut    int
	volLastChngd time.Time
	muteButton   *widget.Button
	muted        bool
}

func NewVolumeSlider() *VolumeSlider {
	vs := &VolumeSlider{}
	vs.muted = false
	vs.vol = binding.NewFloat()
	vs.volSlider = widget.NewSliderWithData(0, 100, vs.vol)
	vs.volSlider.Orientation = widget.Orientation(f2.OrientationVerticalUpsideDown)
	vs.volSlider.Move(f2.NewPos(0, 20))

	vs.LastVol = State.Volume
	vs.volLastChngd = time.Now()

	vs.volSlider.OnChanged = func(val float64) {
		if time.Since(vs.volLastChngd).Milliseconds() > 100 {
			vs.volLastChngd = time.Now()
			v := int(val)
			State.Volume = state.TrackVolume(v)
			Hw.Request("mpd", fmt.Sprintf("setvol %d", v))
		}
	}
	vs.muteButton = widget.NewButtonWithIcon("", theme.VolumeMuteIcon(), vs.muteTapped)

	return vs
}

func (vs *VolumeSlider) muteTapped() {
	if vs.muted {
		Hw.Request("mpd", fmt.Sprintf("setvol %d", vs.volPremut))
		vs.muteButton.SetIcon(theme.VolumeMuteIcon())
	} else {
		vs.volPremut = int(State.Volume)
		Hw.Request("mpd", "setvol 15")
		vs.muteButton.SetIcon(theme.CancelIcon())
	}

	vs.volSlider.SetValue(float64(State.Volume))

	vs.muted = !vs.muted
}

func (vs *VolumeSlider) getGUI() *f2.Container {
	return container.NewBorder(nil, vs.muteButton, nil, nil, vs.volSlider)
}

func (vs *VolumeSlider) UpdateVolume(v state.TrackVolume) {
	vs.vol.Set(float64(v))
	vs.LastVol = v
	vs.volSlider.SetValue(float64(v))
}

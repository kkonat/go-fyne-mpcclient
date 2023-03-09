package gui

import (
	"fmt"
	fe "remotecc/cmd/fynextensions"
	"remotecc/cmd/state"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type PlayerTab struct {
	// gui components
	artist, album, track *canvas.Text
	prgrs                *fe.TappableProgressBar
	timeButtonn          *widget.Button
	// bindables
	elapsed  binding.Float
	timemode bool
}

func NewPlayerTab() *PlayerTab {
	return &PlayerTab{}
}

func (p *PlayerTab) getGUI() *fyne.Container {
	p.artist = fe.NewText("Artist:", 12)
	p.album = fe.NewText("Album:", 12)
	p.track = fe.NewText("Trk:", 12)

	p.elapsed = binding.NewFloat()
	p.prgrs = fe.NewTappableProgressBarWithData(p.elapsed)
	p.prgrs.OnTapped = p.prgrsBarTap
	p.prgrs.Max = float64(State.Track.Duration)
	p.prgrs.TextFormatter = func() string {
		if p.timemode {
			return "-" + state.TrkTimeToString(float32(float64(State.Track.Duration)-p.prgrs.Value))
		} else {
			return state.TrkTimeToString(float32(p.prgrs.Value))
		}

	}
	p.timeButtonn = widget.NewButtonWithIcon("", theme.LogoutIcon(), p.tapTimeButton)
	tabPlayer := container.NewMax(
		container.NewVBox(
			p.artist, p.album, p.track,
			container.NewBorder(nil, nil, nil, p.timeButtonn, p.prgrs),
			container.NewGridWithColumns(5,
				widget.NewButtonWithIcon("", theme.MediaSkipPreviousIcon(), func() {
					Hw.Request("mpd", "previous")
					MW.StatusLine.Set("skip back", 1)
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
					MW.StatusLine.Set("skip next", 1)
				})),
			// widget.NewSeparator(),
		))
	return tabPlayer
}
func (p *PlayerTab) tapTimeButton() {
	p.timemode = !p.timemode
	var icon fyne.Resource
	if p.timemode {
		icon = theme.LoginIcon()
	} else {
		icon = theme.LogoutIcon()
	}
	p.timeButtonn.SetIcon(icon)
}
func (p *PlayerTab) prgrsBarTap(percentPos float64) {
	seekTo := float64(State.Track.Duration) * percentPos
	song := State.Track.Song
	Hw.Request("mpd", fmt.Sprintf("seek %d %d", song, int(seekTo)))
}

func (p *PlayerTab) UpdateTrackDetails(ti *state.TrackInfo) {
	p.album.Text = "Album: " + ti.Album
	p.artist.Text = "Artist: " + ti.Artist
	p.track.Text = ti.Track
	p.album.Refresh()
	p.artist.Refresh()
	p.track.Refresh()
	p.prgrs.Max = float64(ti.Duration)
}

func (p *PlayerTab) UpdateTrackElapsedTime(elTime state.TrackTime) {
	p.elapsed.Set(float64(elTime))
}

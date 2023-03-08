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
	artist, album, track *canvas.Text
	muteButton           *widget.Button
	slider               *widget.Slider
	vol                  binding.Float
	elapsed              binding.Float
	lastVol              state.TrackVolume
	volLastChngd         time.Time
	prgrs                *fe.TappableProgressBar
	IPaddrs              binding.String
	volPremut            int
	muted                bool
	stateStream          chan any

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
		vol:     binding.NewFloat(),
		elapsed: binding.NewFloat(),
		IPaddrs: binding.NewString(),
		muted:   false,
	}

	initStatusText()

	fe.NewText("Volume:", 10)
	c.IPaddrs.Set("192.168.0.95:6600")

	c.slider = widget.NewSliderWithData(0, 100, c.vol)
	c.slider.Orientation = widget.Orientation(f2.OrientationVerticalUpsideDown)
	c.slider.Move(f2.NewPos(0, 20))

	c.lastVol = State.Volume
	c.volLastChngd = time.Now()

	// volume slider
	// slider dragging
	c.slider.OnChanged = func(val float64) {
		if time.Since(c.volLastChngd).Milliseconds() > 100 {
			c.volLastChngd = time.Now()
			v := int(val)
			State.Volume = state.TrackVolume(v)
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

	tabPlaylist := getTabPlaylist(Hw)

	tabPlayer := container.NewMax(
		container.NewVBox(
			c.artist, c.album, c.track, c.prgrs,
			container.NewGridWithColumns(5,
				widget.NewButtonWithIcon("", theme.MediaSkipPreviousIcon(), func() {
					Hw.Request("mpd", "previous")
					setStatusText("skip back", 1)
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
					setStatusText("skip next", 1)
				})),
			// widget.NewSeparator(),
			widget.NewLabel("Playlist:"),
			tabPlaylist),
	)

	tabHW := getTabHW(Hw, State)

	tabSettings :=
		container.NewVBox(
			widget.NewLabel("Settings"),
			widget.NewSeparator(),
			widget.NewLabel("IP:"),
			widget.NewEntryWithData(c.IPaddrs))

	treeData := map[string][]string{
		"":          {"2015", "2016", "Playlisty"},
		"2015":      {"Album 1", "Album 2"},
		"2016":      {"Album 1", "Album 2"},
		"Playlisty": {"Khruangbin Vibes", "Spotify PlOtW", "blabla", "blabla vibes", "blabla sounds"},
		"Album 1":   {"Track 1", "Track 2"},
	}
	tabFilesTree := widget.NewTreeWithStrings(treeData)

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("", theme.MediaPlayIcon(), tabPlayer),
		container.NewTabItemWithIcon("", theme.StorageIcon(), tabPlaylist),
		container.NewTabItemWithIcon("", theme.MediaMusicIcon(), tabFilesTree),
		container.NewTabItemWithIcon("", theme.ComputerIcon(), tabHW),
		container.NewTabItemWithIcon("", theme.SettingsIcon(), tabSettings),
	)
	tabs.SetTabLocation(container.TabLocationTop)

	c.muteButton = widget.NewButtonWithIcon("", theme.VolumeMuteIcon(), func() {
		if c.muted { //unmute
			Hw.Request("mpd", fmt.Sprintf("setvol %d", c.volPremut))
			c.muteButton.SetIcon(theme.VolumeMuteIcon())

		} else {
			c.volPremut = int(c.State.Volume)
			Hw.Request("mpd", "setvol 15")
			c.muteButton.SetIcon(theme.CancelIcon())
		}
		c.slider.SetValue(float64(c.State.Volume))
		c.muted = !c.muted
	})

	w.SetContent(
		container.NewBorder(
			nil,
			nil,
			nil,
			container.NewBorder(nil, c.muteButton, nil, nil, c.slider),                                      // R
			container.NewBorder(nil, container.NewVBox(widget.NewSeparator(), statusLine), nil, nil, tabs))) // center

	return c
}

func (c ControlCenterPanelGUI) UpdateOnlineStatus(online bool, ps state.PlayStatus) {
	if !online {
		setStatusText("Offline. Waiting for connection...", 0)
	} else {
		c.UpdatePlayStatus(ps)
	}
}

func (c *ControlCenterPanelGUI) UpdatePlayStatus(s state.PlayStatus) {
	switch s {
	case state.Playing:
		setStatusText("Playing...", 0)
	case state.Stopped:
		setStatusText("Stopped", 0)
	case state.Paused:
		setStatusText("Paused", 0)
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
	c.slider.SetValue(float64(v))
}

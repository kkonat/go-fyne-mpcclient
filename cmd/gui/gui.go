package gui

import (
	hw "remotecc/cmd/hwinterface"
	"remotecc/cmd/state"

	f2 "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

type MainWindow struct {

	// tabs
	PlaylistTab  *PlaylistTab
	PlayerTab    *PlayerTab
	HWTab        *HWTab
	SettingsTab  *SettingsTab
	FilesTreeTab *FilesTreeTab
	StatusLine   *StatusText

	// complex elements
	VolSlider *VolumeSlider
}

// global references for all GUIC components
// declared global so they do not need to be passed to each GUI element
var (
	Hw          *hw.HWInterface
	State       *state.PlayerState
	stateStream chan any
	MW          *MainWindow
)

func NewMainWindow() *MainWindow {
	MW = &MainWindow{}

	return MW
}

func Init(w f2.Window, stream chan any, s *state.PlayerState, h *hw.HWInterface) {

	// Set global variables
	State = s
	Hw = h
	stateStream = stream

	MW = NewMainWindow()

	// Create complex controls
	MW.StatusLine = NewStatusText()
	MW.VolSlider = NewVolumeSlider()

	// Create Tabs
	MW.PlaylistTab = NewPlaylistTab()
	MW.PlayerTab = NewPlayerTab()
	MW.HWTab = NewHWTab()
	MW.SettingsTab = NewSettingsTab()
	MW.FilesTreeTab = NewFilesTreeTab()

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("", theme.MediaPlayIcon(), MW.PlayerTab.getGUI()),
		container.NewTabItemWithIcon("", theme.StorageIcon(), MW.PlaylistTab.getGUI()),
		container.NewTabItemWithIcon("", theme.MediaMusicIcon(), MW.FilesTreeTab.getGUI()),
		container.NewTabItemWithIcon("", theme.ComputerIcon(), MW.HWTab.getGUI()),
		container.NewTabItemWithIcon("", theme.SettingsIcon(), MW.SettingsTab.getGUI()),
	)
	tabs.SetTabLocation(container.TabLocationTop)

	w.SetContent(
		container.NewBorder(
			nil,
			nil,
			nil,
			MW.VolSlider.getGUI(),   // R
			container.NewBorder(nil, MW.StatusLine.getGUI(), nil, nil, tabs))) // center

}

func (mw MainWindow) UpdateOnlineStatus(online bool, ps state.PlayStatus) {
	if !online {
		mw.StatusLine.Set("Offline. Waiting for connection...", 0)
	} else {
		mw.UpdatePlayStatus(ps)
	}
}

func (mw *MainWindow) UpdatePlayStatus(s state.PlayStatus) {
	switch s {
	case state.Playing:
		mw.StatusLine.Set("Playing...", 0)
	case state.Stopped:
		mw.StatusLine.Set("Stopped", 0)
	case state.Paused:
		mw.StatusLine.Set("Paused", 0)
	}
}



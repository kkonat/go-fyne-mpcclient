package gui

import (
	hw "remotecc/cmd/hwinterface"
	"remotecc/cmd/state"
	"remotecc/cmd/storage"

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
	tabs         *container.AppTabs

	// complex elements
	VolSlider *VolumeSlider
}

// references to all GUI components
// declared global so they do not need to be passed to each GUI element
var (
	Hw          *hw.HWInterface
	State       *state.PlayerState
	stateStream chan any
	MW          *MainWindow
	AppWindow   *f2.Window
)

func NewMainWindow() *MainWindow {
	MW = &MainWindow{}
	return MW
}

func Init(w *f2.Window, stream chan any, s *state.PlayerState, h *hw.HWInterface) {
	// Set global variables
	AppWindow = w
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

	MW.regenWindow()

}
func (mw *MainWindow) regenWindow() {
	mw.tabs = container.NewAppTabs()
	mw.tabs.Append(container.NewTabItemWithIcon("", theme.MediaPlayIcon(), MW.PlayerTab.getGUI()))
	mw.tabs.Append(container.NewTabItemWithIcon("", theme.StorageIcon(), MW.PlaylistTab.getGUI()))
	mw.tabs.Append(container.NewTabItemWithIcon("", theme.MediaMusicIcon(), MW.FilesTreeTab.getGUI()))
	if storage.AppSettings.ShowHWCtrl {
		mw.tabs.Append(container.NewTabItemWithIcon("", theme.ComputerIcon(), MW.HWTab.getGUI()))
	}
	mw.tabs.Append(container.NewTabItemWithIcon("", theme.SettingsIcon(), MW.SettingsTab.getGUI()))
	mw.tabs.SetTabLocation(container.TabLocationTop)
	(*AppWindow).SetContent(
		container.NewBorder(
			nil,
			nil,
			nil,
			MW.VolSlider.getGUI(), // R
			container.NewBorder(nil, MW.StatusLine.getGUI(), nil, nil, MW.tabs))) // center
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

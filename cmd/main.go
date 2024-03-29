package main

import (
	"context"
	"log"
	"remotecc/cmd/gui"
	hw "remotecc/cmd/hwinterface"
	"remotecc/cmd/state"
	"remotecc/cmd/storage"
	"time"

	f2 "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
)

var (
	stateStream = make(chan any)
	ctx         context.Context
	Hw          *hw.HWInterface
	State       *state.PlayerState
	_           f2.Theme = (*myTheme)(nil)
)

func main() {
	contx, cancel := context.WithCancel(context.Background())
	ctx = contx
	defer cancel()

	storage.Init()
	Hw = hw.NewHWInterface()
	State = state.NewPlayerState(Hw, stateStream)

	fyneApp := app.New()
	window := fyneApp.NewWindow("Remote Control Center")
	fyneApp.Settings().SetTheme(&myTheme{})
	window.Resize(f2.Size{Width: 300, Height: 500})

	gui.Init(&window, stateStream, State, Hw)

	go MonitorStateChanges()
	go HandleStateChanges()

	window.ShowAndRun()

	// storage.SaveData()
	log.Print("Goodbye")
}

func MonitorStateChanges() {

	updTrackData := time.Tick(time.Millisecond * 500)
	updPlayerStatus := time.Tick(time.Millisecond * 85)
	updHwStatus := time.Tick(time.Second)

	for {
		select {

		case <-updTrackData:
			go State.GetTrackData()

		case <-updPlayerStatus:
			go State.GetMPDStatus()

		case <-updHwStatus:
			go State.GetHWState()

		case <-ctx.Done():
			return
		}
	}
}

func HandleStateChanges() {
	for {
		select {
		case chgdState := <-stateStream:
			switch newValue := chgdState.(type) {

			case state.TrackInfo:
				gui.MW.PlayerTab.UpdateTrackDetails(&newValue)
				canvas.Refresh((*gui.AppWindow).Canvas().Content())

			case state.TrackVolume:
				gui.MW.VolSlider.UpdateVolume(newValue)

			case state.PlayStatus:
				gui.MW.UpdatePlayStatus(state.PlayStatus(newValue))

			case state.TrackTime:
				gui.MW.PlayerTab.UpdateTrackElapsedTime(newValue)

			case state.PowerStatus:
				gui.MW.HWTab.UpdatePowerBbutton(bool(newValue))

			case state.OnlineStatus:
				gui.MW.UpdateOnlineStatus(bool(newValue), State.Status)
			}

		case <-ctx.Done():
			return
		}
	}
}

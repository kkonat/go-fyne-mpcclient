package main

import (
	"log"
	"time"

	"remotecc/cmd/gui"
	hw "remotecc/cmd/hwinterface"
	"remotecc/cmd/state"

	"fyne.io/fyne/v2/app"
	"golang.org/x/net/context"
)

var (
	buildDate   string
	stateStream = make(chan any)
	ctx         context.Context
	Hw          *hw.HWInterface    = hw.NewHWInterface()
	State       *state.PlayerState = state.NewPlayerState(Hw, stateStream)
	ccgui       *gui.ControlCenterPanelGUI
)

func main() {
	// fmt.Println("Build date: ", buildDate)
	contx, cancel := context.WithCancel(context.Background())
	ctx = contx
	defer cancel()

	fyneApp := app.New()
	window := fyneApp.NewWindow("Remote Control Center")

	ccgui = gui.New(window, stateStream, State, Hw)

	go MonitorStateChanges()
	go HandleStateChanges()

	window.ShowAndRun()

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
				ccgui.UpdateTrackDetails(&newValue)

			case state.TrackVolume:
				ccgui.UpdateVolume(newValue)

			case state.PlayStatus:
				ccgui.UpdatePlayStatus(state.PlayStatus(newValue))

			case state.TrackTime:
				ccgui.UpdateTrackElapsedTime(newValue)

			case state.PowerStatus:
				ccgui.UpdatePowerBbutton(bool(newValue))

			case state.OnlineStatus:
				ccgui.UpdateOnlineStatus(bool(newValue), State.Status)
			}

		case <-ctx.Done():
			return
		}
	}
}

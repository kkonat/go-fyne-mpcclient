package main

import (
	"context"
	"log"
	"time"
)

type ControlCenterApp struct {
	ctx         context.Context
	state       *PlayerState
	hw          *HWInterface
	stateStream chan any
}

func NewControlCenterApp(ctx context.Context, hw *HWInterface) *ControlCenterApp {
	app := &ControlCenterApp{
		ctx: ctx,
		hw:  hw,
	}
	app.stateStream = make(chan any)
	app.state = NewPlayerState(hw, &app.stateStream)

	return app
}

func (a *ControlCenterApp) toglleHWPower() (bool, error) {
	a.hw.togglePower()
	return a.state.getHWState()
}
func (a *ControlCenterApp) refreshState() {

	updTrackData := time.Tick(time.Millisecond * 500)
	updPlayerStatus := time.Tick(time.Millisecond * 85)
	updHwStatus := time.Tick(time.Second)
loop:
	for {
		select {

		case <-updTrackData:
			go a.state.getTrkData()

		case <-updPlayerStatus:
			go a.state.getStatus()

		case <-updHwStatus:
			go a.state.getHWState()

		case <-a.ctx.Done():
			break loop
		}
	}
}

func (a *ControlCenterApp) updateGUI(gui *ControlCenterPanelGUI) {

loop:
	for {
		select {
		case what := <-a.stateStream:
			switch newValue := what.(type) {

			case TrackInfo:
				//	log.Printf("update track info: %v", newValue)
				gui.updateTrackDetails(&newValue)

			case TrackVolume:
				gui.updateVolume(newValue)

			case PlayStatus:
				log.Printf("updae play status: %v", newValue)
				// TODO update player status

			case TrackTime:
				gui.updateTrackElapsedTime(newValue)

			case PowerStatus:
				log.Printf("power status: %v", newValue)
				gui.updatePowerBbutton(bool(newValue))
			}

		case <-a.ctx.Done():
			break loop
		}
	}
}

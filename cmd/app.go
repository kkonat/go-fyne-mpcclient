package main

import (
	"context"
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

func (a *ControlCenterApp) monitorStateChanges() {

	updTrackData := time.Tick(time.Millisecond * 500)
	updPlayerStatus := time.Tick(time.Millisecond * 85)
	updHwStatus := time.Tick(time.Second)

	for {
		select {

		case <-updTrackData:
			go a.state.getTrackData()

		case <-updPlayerStatus:
			go a.state.getMPDStatus()

		case <-updHwStatus:
			go a.state.getHWState()

		case <-a.ctx.Done():
			return
		}
	}
}

func (a *ControlCenterApp) handleStateChanges(gui *ControlCenterPanelGUI) {
	for {
		select {
		case chgdState := <-a.stateStream:
			switch newValue := chgdState.(type) {

			case TrackInfo:
				//	log.Printf("update track info: %v", newValue)
				gui.updateTrackDetails(&newValue)

			case TrackVolume:
				gui.updateVolume(newValue)

			case PlayStatus:
				gui.UpdatePlayStatus(PlayStatus(newValue))

			case TrackTime:
				gui.updateTrackElapsedTime(newValue)

			case PowerStatus:
				// log.Printf("power status: %v", newValue)
				gui.updatePowerBbutton(bool(newValue))
			case OnlineStatus:

				gui.updateOnlineStatus(bool(newValue), a.state.status)
			}

		case <-a.ctx.Done():
			return
		}
	}
}

package main

import (
	"context"
	"log"
	"time"
)

type ControlCenterApp struct {
	mpdClient   *Client
	ctrlClient  *Client
	ctx         context.Context
	state       *playerState
	stateChange chan any
}

func NewControlCenterApp(ctx context.Context) *ControlCenterApp {
	app := &ControlCenterApp{
		mpdClient:  NewClient("192.168.0.95:6600"),
		ctrlClient: NewClient("192.168.0.95:1025"),
		ctx:        ctx,
	}
	app.stateChange = make(chan any)
	app.state = NewState(app, &app.stateChange)

	return app
}

func (a *ControlCenterApp) refreshState() {
	
	updTrackData := time.Tick(time.Millisecond * 500)
	updPlayerStatus := time.Tick(time.Millisecond * 100)

	loop:
	for {
		select {

		case <-updTrackData:
			go a.state.getTrkData()

		case <-updPlayerStatus:
			go a.state.getStatus()

		case <-a.ctx.Done():
			break loop
		}
	}
}

func (a *ControlCenterApp) updateGUI(gui *ControlCenterPanelGUI) {

loop:
	for {
		select {
		case what := <-a.stateChange:
			switch newValue := what.(type) {

			case TrackInfo:
				log.Printf("update track info: %v", newValue)
				gui.updateTrackDetails(&newValue)

			case TrackVolume:
				log.Printf("updae track volume: %v", newValue)
				gui.updateVolume(newValue)

			case PlayStatus:
				log.Printf("updae play status: %v", newValue)
				// TODO update player status

			case TrackTime:
				// log.Printf("updae track time elapsed: %v", newValue)
				gui.updateTrackElapsedTime(newValue)
			}

		case <-a.ctx.Done():
			break loop
		}
	}
}

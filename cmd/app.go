package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
)

type ControlCenterApp struct {
	mpdClient   *Client
	ctrlClient  *Client
	ctx         context.Context
	state       *PlayerState
	stateStream chan any
}

func NewControlCenterApp(ctx context.Context) *ControlCenterApp {
	app := &ControlCenterApp{
		mpdClient:  NewClient("192.168.0.95:6600", false),
		ctrlClient: NewClient("192.168.0.95:1025", true),
		ctx:        ctx,
	}
	app.stateStream = make(chan any)
	app.state = NewPlayerState(app, &app.stateStream)

	return app
}

func (a *ControlCenterApp) refreshState() {

	updTrackData := time.Tick(time.Millisecond * 500)
	updPlayerStatus := time.Tick(time.Millisecond * 85)

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
			}

		case <-a.ctx.Done():
			break loop
		}
	}
}

func (a *ControlCenterApp) chkPowerStatus() (bool, error) {
	res, err := a.ctrlClient.Request("check_extpower")
	fmt.Println(res)
	if err == nil {
		return strings.Split(res[0], ": ")[1] == "1", nil
	}
	return false, err
}

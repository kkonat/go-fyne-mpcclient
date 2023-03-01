package main

import (
	"time"

	"fyne.io/fyne/v2/app"
	"golang.org/x/net/context"
)

const ctrlSrv = "192.168.0.95:1025"
const mpdSrv = "192.168.0.95:6600"

func startMPCmonitor(ctx context.Context, updateCh chan any, ci *playerState) {

	ticker := time.Tick(time.Millisecond * 300)

	mpc := MpcClient{serverAddr: mpdSrv, ctx: ctx, updateCh: updateCh}

loop:
	for {
		select {
		case <-ticker:
			mpc.Update()
		case <-ctx.Done():
			break loop
		}
	}
}

func main() {

	a := app.New()
	w := a.NewWindow("Remote Control Center")

	playerState := playerState{}
	ccp := newControlCenterPanel(w, &playerState)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	updateStream := make(chan any)
	startMPCmonitor(ctx, updateStream, &playerState)

loop:

	for {
		select {
		case what := <-updateStream:
			switch newValue := what.(type) {
			case TrackInfo:
				ccp.updateTrackDetails(newValue)
			case TrackVolume:
				ccp.updateVolume(newValue)
			case PlayerStatus:
				// TODO update player status			case TrackTime:
				// TODO update elaapsed scroller
			}

		case <-ctx.Done():
			break loop
		}
	}

	w.ShowAndRun()
}

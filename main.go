package main

import (
	"log"
	"time"

	"fyne.io/fyne/v2/app"
	"golang.org/x/net/context"
)

func startMPCmonitor(mpcc *MpcClient, ccp *ControlCenterPanel) {

	ticker := time.Tick(time.Millisecond * 300)

loop:
	for {
		select {
		case what := <-mpcc.updateStream:
			switch newValue := what.(type) {
			case TrackInfo:
				log.Print("updae track info")
				ccp.updateTrackDetails(newValue)
			case TrackVolume:
				log.Print("updae track volume")
				ccp.updateVolume(newValue)
			case PlayStatus:
				log.Print("updae play status")
				// TODO update player status
			case TrackTime:
				log.Print("updae track time elapsed")
				// TODO update elaapsed scroller
			}
		case <-ticker:
			log.Print("tick update")
			mpcc.Update()
		case <-mpcc.ctx.Done():
			break loop
		}
	}
}

const ctrlSrvAddr = "192.168.0.95:1025"
const mpdSrvAddr = "192.168.0.95:6600"

func main() {

	a := app.New()
	window := a.NewWindow("Remote Control Center")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	updateStream := make(chan any, 10)
	mpc := &MpcClient{
		mpdServerAddr: mpdSrvAddr,
		ctrlSrvAddr:   ctrlSrvAddr,
		ctx:           ctx,
		updateStream:  updateStream,
	}
	ccp := newControlCenterPanelGUI(window, mpc)
	log.Print("Starting MPC monitor")
	go startMPCmonitor(mpc, ccp)
	log.Print("Showing window")
	window.ShowAndRun()

	log.Print("Main loop end")
}

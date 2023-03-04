package main

import (
	"log"

	"fyne.io/fyne/v2/app"
	"golang.org/x/net/context"
)

func main() {
	fyneApp := app.New()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	updateStream := make(chan any, 10)

	ccApp := NewControlCenterApp(ctx, updateStream)

	window := fyneApp.NewWindow("Remote Control Center")
	ccGui := newControlCenterPanelGUI(window, ccApp)

	go ccApp.state.monitor(ccGui)

	window.ShowAndRun()

	log.Print("Goodbye")
}

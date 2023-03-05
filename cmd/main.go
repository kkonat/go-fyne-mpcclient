package main

import (
	"log"

	"fyne.io/fyne/v2/app"
	"golang.org/x/net/context"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fyneApp := app.New()
	hw := NewHWInterface()
	ccApp := NewControlCenterApp(ctx, hw)
	window := fyneApp.NewWindow("Remote Control Center")
	ccGui := newControlCenterAppGUI(window, ccApp)

	go ccApp.refreshState()
	go ccApp.updateGUI(ccGui)

	window.ShowAndRun()

	log.Print("Goodbye")
}

package main

import (
	"log"

	"fyne.io/fyne/v2/app"
	"golang.org/x/net/context"
)

func main() {

	a := app.New()
	window := a.NewWindow("Remote Control Center")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	updateStream := make(chan any, 10)

	ccApp := NewControlCenterApp(ctx, updateStream)

	ccGui := newControlCenterPanelGUI(window, ccApp)

	go ccApp.monitorState(ccGui)

	window.ShowAndRun()

	log.Print("Goodbye")
}

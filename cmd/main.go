package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2/app"
	"golang.org/x/net/context"
)

var buildDate string

func main() {
	fmt.Println("Build date: ", buildDate)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fyneApp := app.New()
	hw := NewHWInterface()
	ccApp := NewControlCenterApp(ctx, hw)
	window := fyneApp.NewWindow("Remote Control Center")
	ccGui := newControlCenterAppGUI(window, ccApp)

	go ccApp.monitorStateChanges()
	go ccApp.handleStateChanges(ccGui)

	window.ShowAndRun()

	log.Print("Goodbye")
}

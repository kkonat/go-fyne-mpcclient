package gui

import (
	fe "remotecc/cmd/fynextensions"
	"time"

	"fyne.io/fyne/v2/canvas"
)

var statusLine *canvas.Text
var storedStatusText string

func initStatusText() {
	statusLine = fe.NewText("Ready", 10)
}
func setStatusText(newText string, howLong int) {
	if howLong == 0 {
		statusLine.Text = newText
		statusLine.Refresh()
		return
	}
	storedStatusText = statusLine.Text
	statusLine.Text = newText
	time.AfterFunc(time.Second*time.Duration(howLong),
		func() {
			statusLine.Text = storedStatusText
			statusLine.Refresh()
		})
	statusLine.Refresh()
}

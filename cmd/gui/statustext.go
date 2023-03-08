package gui

import (
	fe "remotecc/cmd/fynextensions"
	"time"

	f2 "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type StatusText struct {
	statusLine       *canvas.Text
	storedStatusText string
}

func NewStatusText() *StatusText {
	st := &StatusText{}
	st.statusLine = fe.NewText("Ready", 10)

	return st
}

func (st *StatusText) Set(newText string, howLong int) {
	if howLong == 0 {
		st.statusLine.Text = newText
		st.statusLine.Refresh()
		return
	}
	st.storedStatusText = st.statusLine.Text
	st.statusLine.Text = newText
	time.AfterFunc(
		time.Second*time.Duration(howLong),
		st.restorePreviousText)
	st.statusLine.Refresh()
}

func (st *StatusText) restorePreviousText() {
	st.statusLine.Text = st.storedStatusText
	st.statusLine.Refresh()
}

func (st *StatusText) getGUI() *f2.Container {
	return container.NewVBox(widget.NewSeparator(), st.statusLine)
}

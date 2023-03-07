package fynextensions

import (
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)

func NewText(text string, size float32) *canvas.Text {
	fg := theme.ForegroundColor()
	return &canvas.Text{
		Color:    fg,
		Text:     text,
		TextSize: size,
	}
}

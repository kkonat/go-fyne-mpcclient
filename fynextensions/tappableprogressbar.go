package fynextensions

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type TappableProgressBar struct {
	widget.ProgressBar
	OnTapped func(newVal float64)
}

func NewTappableProgressBar() *TappableProgressBar {
	tpb := &TappableProgressBar{}
	tpb.ExtendBaseWidget(tpb)

	return tpb
}

func NewTappableProgressBarWithData(data binding.Float) *TappableProgressBar {
	p := NewTappableProgressBar()
	p.ExtendBaseWidget(p)
	p.Bind(data)
	return p
}

func (t *TappableProgressBar) toPercent(pe *fyne.PointEvent) float64 {
	// log.Printf("Tapped, pe=%+v\n", pe)
	width := t.Size().Width
	tp := pe.Position.X
	return float64(tp) / float64(width)
}

func (t *TappableProgressBar) Tapped(pe *fyne.PointEvent) {
	t.OnTapped(t.toPercent(pe))
}

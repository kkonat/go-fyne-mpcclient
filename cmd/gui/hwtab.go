package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type HWTab struct {
	bInputPlayer, bInputTv, bPower, bShtDn *widget.Button
}

func NewHWTab() *HWTab {
	t := &HWTab{}
	t.bInputPlayer = widget.NewButton("Player", func() {
		MW.StatusLine.Set("switched to Player", 3)
		Hw.Request("ctrl", "deq_input_coaxial")
	})

	t.bInputTv = widget.NewButton("TV", func() {
		MW.StatusLine.Set("switched to TV", 3)
		Hw.Request("ctrl", "deq_input_optical")
	})

	t.bShtDn = widget.NewButton("Shutdown", func() {
		MW.StatusLine.Set("Shutting down...", 0)
		Hw.Request("ctrl", "server_poweroff")
	})
	t.bPower = widget.NewButton("Power", func() {
		t.TogglePower()
	})
	return t
}

func (t *HWTab) getGUI() *fyne.Container {

	return container.NewVBox(
		container.NewGridWithColumns(2,
			t.bInputPlayer,
			t.bInputTv),
		t.bPower,
		t.bShtDn)
}

func (t *HWTab) UpdatePowerBbutton(powerOn bool) {
	if powerOn {
		t.bPower.SetText("Power off")
	} else {
		t.bPower.SetText("Power on")
	}
}

func (t *HWTab) TogglePower() {
	var state bool
	var err error
	if err = Hw.TogglePower(); err == nil {
		if state, err = State.GetHWState(); err == nil {
			if state {
				MW.StatusLine.Set("powered off", 3)
			} else {
				MW.StatusLine.Set("powered on", 3)
			}
			t.UpdatePowerBbutton(state)
		}
	}
}

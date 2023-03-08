package gui

import (
	"remotecc/cmd/hwinterface"

	"remotecc/cmd/state"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var Hw *hwinterface.HWInterface
var State *state.PlayerState

var bInputPlayer = widget.NewButton("Player", func() {
	setStatusText("switched to Player", 3)
	Hw.Request("ctrl", "deq_input_coaxial")
})

var bInputTV = widget.NewButton("TV", func() {
	setStatusText("switched to TV", 3)
	Hw.Request("ctrl", "deq_input_optical")
})

var bPower *widget.Button

var bShtDn = widget.NewButton("Shutdown", func() {
	setStatusText("Shutting down...", 0)
	Hw.Request("ctrl", "server_poweroff")
})

func UpdatePowerBbutton(powerOn bool) {
	if powerOn {
		bPower.SetText("Power off")
	} else {
		bPower.SetText("Power on")
	}
}

func TogglePower() {
	var state bool
	var err error
	if err = Hw.TogglePower(); err == nil {
		if state, err = State.GetHWState(); err == nil {
			if state {
				setStatusText("powered off", 3)
			} else {
				setStatusText("powered on", 3)
			}
			UpdatePowerBbutton(state)
		}
	}
}

func getTabHW(hw *hwinterface.HWInterface, state *state.PlayerState) *fyne.Container {
	bPower = widget.NewButton("Power", func() {
		TogglePower()
	})
	Hw = hw
	State = state
	return container.NewVBox(
		container.NewGridWithColumns(2,
			bInputPlayer,
			bInputTV),
		bPower,
		bShtDn)
}

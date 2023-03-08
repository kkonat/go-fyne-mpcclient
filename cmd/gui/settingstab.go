package gui

import (
	"remotecc/cmd/storage"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type SettingsTab struct {
	IPaddrs binding.String
}

func NewSettingsTab() *SettingsTab {
	t := &SettingsTab{}

	t.IPaddrs = binding.NewString()
	t.IPaddrs.Set(storage.AppSettings.Server.IPAddr + ":" + storage.AppSettings.Server.MPDPort)
	return t
}

func (t *SettingsTab) dfltValEntry(txt string) *widget.Entry {
	w := widget.NewEntry()
	w.SetText(txt)
	return w
}

func (t *SettingsTab) getGUI() *widget.Form {
	form := widget.NewForm(
		widget.NewFormItem("IP address:", t.dfltValEntry(storage.AppSettings.Server.IPAddr)),
		widget.NewFormItem("MPD port:", t.dfltValEntry(storage.AppSettings.Server.MPDPort)),
		widget.NewFormItem("HW ctrl port:", t.dfltValEntry(storage.AppSettings.Server.CtrlPort)),
	)
	return form
}

package gui

import (
	"errors"
	"remotecc/cmd/storage"
	"strconv"
	"strings"

	f2 "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type SettingsTab struct {
	form        *widget.Form
	showHWctrls *widget.Check
}

func NewSettingsTab() *SettingsTab {
	t := &SettingsTab{}
	t.showHWctrls = widget.NewCheck("Show HW controls", func(bool) {})
	return t
}

func (t *SettingsTab) dfltValEntry(txt string, chk f2.StringValidator) *widget.Entry {
	w := widget.NewEntry()
	w.SetText(txt)
	// w.SetPlaceHolder(txt)
	w.Validator = chk
	return w
}

func (t *SettingsTab) getGUI() *f2.Container {
	t.showHWctrls.SetChecked(storage.AppSettings.ShowHWCtrl)
	t.showHWctrls.OnChanged = t.chngHWctrls
	t.form = widget.NewForm(
		widget.NewFormItem("IP address:", t.dfltValEntry(storage.AppSettings.Server.IPAddr, isIPaddr)),
		widget.NewFormItem("MPD port:", t.dfltValEntry(storage.AppSettings.Server.MPDPort, isWord)),
		widget.NewFormItem("HW ctrl port:", t.dfltValEntry(storage.AppSettings.Server.CtrlPort, isWord)),
	)
	return container.NewVBox(
		t.form,
		t.showHWctrls,
		widget.NewButton("Save...", func() {
			err := t.form.Validate()
			if err == nil {
				dialog.ShowConfirm("Saving settings", "Do you want to save these settings?", func(save bool) {
					if save {
						t.extractAndSave()
					}
				}, (*AppWindow))
			} else {
				dialog.ShowError(errors.New("Invalid data. Can't save"), (*AppWindow))
			}
		}),
	)
}

func (t *SettingsTab) chngHWctrls(newVal bool) {

	storage.AppSettings.ShowHWCtrl = newVal
	MW.regenWindow()
	MW.tabs.SelectIndex(len(MW.tabs.Items) - 1)
}

func (t *SettingsTab) extractAndSave() {
	storage.AppSettings.Server.IPAddr = t.form.Items[0].Widget.(*widget.Entry).Text
	storage.AppSettings.Server.MPDPort = t.form.Items[1].Widget.(*widget.Entry).Text
	storage.AppSettings.Server.CtrlPort = t.form.Items[2].Widget.(*widget.Entry).Text

	(*AppWindow).Content().Refresh()
	storage.SaveData()
}

func isIPaddr(s string) error {
	bytes := strings.Split(s, ".")
	if len(bytes) != 4 {
		return errors.New("invalid")
	}
	for _, b := range bytes {
		if isByte(b) != nil {
			return errors.New("invalid")
		}
	}
	return nil
}
func isWord(s string) error {
	val, err := strconv.ParseUint(s, 10, 32)
	if err == nil && val < 65536 {
		return nil
	}
	return errors.New("invalid")
}
func isByte(s string) error {
	val, err := strconv.ParseUint(s, 10, 32)
	if err == nil && val < 256 {
		return nil
	}
	return errors.New("invalid")
}

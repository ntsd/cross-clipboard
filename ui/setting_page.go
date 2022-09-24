package ui

import (
	"strconv"
	"unicode"

	"github.com/rivo/tview"
)

func numberValidator(text string, lastChar rune) bool {
	if !unicode.IsDigit(lastChar) {
		return false
	}
	return true
}

func (v *View) newSettingPage() *Page {
	form := tview.NewForm().
		// clipboard
		AddCheckbox("encrypt enabled", v.CrossClipboard.Config.EncryptEnabled, nil).
		AddInputField("max size (bytes)", strconv.Itoa(v.CrossClipboard.Config.MaxSize), 10, numberValidator, nil).
		AddInputField("max history", strconv.Itoa(v.CrossClipboard.Config.MaxHistory), 3, numberValidator, nil).

		// device
		AddCheckbox("auto trust device", v.CrossClipboard.Config.AutoTrust, nil).

		// ui
		AddCheckbox("hidden text", v.CrossClipboard.Config.HiddenText, nil).
		AddCheckbox("terminal mode", v.CrossClipboard.Config.TerminalMode, nil).

		// network
		AddInputField("groupname", v.CrossClipboard.Config.GroupName, 50, nil, nil).
		AddInputField("protocal id", v.CrossClipboard.Config.ProtocolID, 50, nil, nil).
		AddInputField("port", strconv.Itoa(v.CrossClipboard.Config.ListenPort), 5, numberValidator, nil).

		// submit
		AddButton("save", func() {
			// TODO
		}).
		AddButton("cancel", func() {
			// TODO
		}).
		AddButton("default", func() {
			// TODO
		})

	form.SetBorder(true).SetTitle("setting").SetTitleAlign(tview.AlignCenter)

	return &Page{
		Title:   "Setting",
		Content: form,
	}
}

package ui

import (
	"log"
	"strconv"
	"unicode"

	"github.com/ntsd/cross-clipboard/pkg/xerror"
	"github.com/rivo/tview"
)

func numberValidator(text string, lastChar rune) bool {
	if !unicode.IsDigit(lastChar) {
		return false
	}
	return true
}

func unsafeStringToInt(text string) int {
	n, err := strconv.Atoi(text)
	if err != nil {
		log.Fatal(err)
	}
	return n
}

func (v *View) newSettingPage() *Page {
	cfg := v.CrossClipboard.Config

	settingLayout := tview.NewGrid()

	settingForm := tview.NewForm()
	settingForm.SetBorder(true).
		SetTitle("Setting").
		SetTitleAlign(tview.AlignCenter)
	settingLayout.AddItem(settingForm, 0, 0, 1, 1, 0, 0, true)

	advanceForm := tview.NewForm()
	advanceForm.SetBorder(true).
		SetTitle("Advance Setting").
		SetTitleAlign(tview.AlignCenter)
	settingLayout.AddItem(advanceForm, 0, 1, 1, 1, 0, 0, false)

	// Add setting form item
	settingForm.AddFormItem(tview.NewCheckbox().
		SetLabel("hidden text").
		SetChecked(v.CrossClipboard.Config.HiddenText).
		SetChangedFunc(func(checked bool) { cfg.HiddenText = checked }))

	settingForm.AddFormItem(tview.NewInputField().
		SetLabel("max size (bytes)").
		SetText(strconv.Itoa(v.CrossClipboard.Config.MaxSize)).
		SetFieldWidth(10).
		SetAcceptanceFunc(numberValidator).
		SetChangedFunc(func(text string) { cfg.MaxSize = unsafeStringToInt(text) }))

	settingForm.AddFormItem(tview.NewInputField().
		SetLabel("max history").
		SetText(strconv.Itoa(v.CrossClipboard.Config.MaxHistory)).
		SetFieldWidth(3).
		SetAcceptanceFunc(numberValidator).
		SetChangedFunc(func(text string) { cfg.MaxHistory = unsafeStringToInt(text) }))

	settingForm.AddFormItem(tview.NewCheckbox().
		SetLabel("auto trust device").
		SetChecked(v.CrossClipboard.Config.AutoTrust).
		SetChangedFunc(func(checked bool) { cfg.AutoTrust = checked }))

	settingForm.AddFormItem(tview.NewCheckbox().
		SetLabel("encrypt enabled").
		SetChecked(v.CrossClipboard.Config.EncryptEnabled).
		SetChangedFunc(func(checked bool) { cfg.EncryptEnabled = checked }))

	// add advance setting form item

	advanceForm.AddFormItem(tview.NewCheckbox().
		SetLabel("terminal mode").
		SetChecked(v.CrossClipboard.Config.TerminalMode).
		SetChangedFunc(func(checked bool) { cfg.TerminalMode = checked }))

	advanceForm.AddFormItem(tview.NewInputField().
		SetLabel("groupname").
		SetText(v.CrossClipboard.Config.GroupName).
		SetFieldWidth(50).
		SetAcceptanceFunc(nil).
		SetChangedFunc(func(text string) { cfg.GroupName = text }))

	advanceForm.AddFormItem(tview.NewInputField().
		SetLabel("host").
		SetText(v.CrossClipboard.Config.ListenHost).
		SetFieldWidth(50).
		SetAcceptanceFunc(nil).
		SetChangedFunc(func(text string) { cfg.ListenHost = text }))

	advanceForm.AddFormItem(tview.NewInputField().
		SetLabel("port").
		SetText(strconv.Itoa(v.CrossClipboard.Config.ListenPort)).
		SetFieldWidth(5).
		SetAcceptanceFunc(numberValidator).
		SetChangedFunc(func(text string) { cfg.ListenPort = unsafeStringToInt(text) }))

	advanceForm.AddButton("save", func() {
		err := cfg.Save()
		if err != nil {
			v.CrossClipboard.ErrorChan <- xerror.NewRuntimeErrorf("can't save config: %v", err)
		}
	})

	advanceForm.AddButton("reset to default", func() {
		// TODO
	})

	return &Page{
		Title:   "Setting",
		Content: settingLayout,
	}
}

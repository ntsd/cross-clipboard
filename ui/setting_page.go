package ui

import (
	"log"
	"os"
	"strconv"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/ntsd/cross-clipboard/pkg/xerror"
	"github.com/rivo/tview"
)

func numberValidator(maxLen int) func(text string, lastChar rune) bool {
	return func(text string, lastChar rune) bool {
		// Note this function will not if it's remove character by back space
		l := len(text)
		if l <= 0 || l >= maxLen {
			return false
		}
		if !unicode.IsDigit(lastChar) {
			return false
		}
		return true
	}
}

func unsafeStringToInt(text string) int {
	if text == "" {
		return 0
	}
	n, err := strconv.Atoi(text)
	if err != nil {
		log.Fatal(err)
	}
	return n
}

func (v *View) newSettingPage() *Page {
	cfg := v.CrossClipboard.Config

	save := func(cb func()) {
		err := cfg.Save()
		if err != nil {
			v.CrossClipboard.ErrorChan <- xerror.NewRuntimeErrorf("can not save config: %v", err)
			return
		}
		if cb != nil {
			cb()
		}
	}

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
	var formItems []tview.FormItem

	formItems = append(formItems, tview.NewCheckbox().
		SetLabel("hidden text").
		SetChecked(v.CrossClipboard.Config.HiddenText).
		SetChangedFunc(func(checked bool) {
			cfg.HiddenText = checked
			save(nil)
		}))

	formItems = append(formItems, tview.NewInputField().
		SetLabel("max size (MB)").
		SetText(strconv.Itoa(v.CrossClipboard.Config.MaxSize)).
		SetFieldWidth(5).
		SetAcceptanceFunc(numberValidator(3)).
		SetChangedFunc(func(text string) {
			cfg.MaxSize = unsafeStringToInt(text)
			save(nil)
		}))

	formItems = append(formItems, tview.NewInputField().
		SetLabel("max history").
		SetText(strconv.Itoa(v.CrossClipboard.Config.MaxHistory)).
		SetFieldWidth(5).
		SetAcceptanceFunc(numberValidator(3)).
		SetChangedFunc(func(text string) {
			cfg.MaxHistory = unsafeStringToInt(text)
			save(nil)
		}))

	formItems = append(formItems, tview.NewCheckbox().
		SetLabel("auto trust device").
		SetChecked(v.CrossClipboard.Config.AutoTrust).
		SetChangedFunc(func(checked bool) {
			cfg.AutoTrust = checked
			save(nil)
		}))

	formItemsLen := len(formItems)
	for _, formItem := range formItems {
		settingForm.AddFormItem(formItem)
	}

	// add advance setting form item
	var advanceFormItems []tview.FormItem

	advanceFormItems = append(advanceFormItems, tview.NewInputField().
		SetLabel("config path").
		SetText(v.CrossClipboard.Config.ConfigDirPath).
		SetFieldWidth(50).
		SetAcceptanceFunc(nil).
		SetChangedFunc(func(text string) {
			cfg.ConfigDirPath = text
		}))

	advanceFormItems = append(advanceFormItems, tview.NewInputField().
		SetLabel("groupname").
		SetText(v.CrossClipboard.Config.GroupName).
		SetFieldWidth(50).
		SetAcceptanceFunc(nil).
		SetChangedFunc(func(text string) {
			cfg.GroupName = text
		}))

	advanceFormItems = append(advanceFormItems, tview.NewInputField().
		SetLabel("host").
		SetText(v.CrossClipboard.Config.ListenHost).
		SetFieldWidth(50).
		SetAcceptanceFunc(nil).
		SetChangedFunc(func(text string) {
			cfg.ListenHost = text
		}))

	advanceFormItems = append(advanceFormItems, tview.NewInputField().
		SetLabel("port").
		SetText(strconv.Itoa(v.CrossClipboard.Config.ListenPort)).
		SetFieldWidth(5).
		SetAcceptanceFunc(numberValidator(5)).
		SetChangedFunc(func(text string) {
			cfg.ListenPort = unsafeStringToInt(text)
		}))

	advanceFormItemsLen := len(advanceFormItems)
	for _, formItem := range advanceFormItems {
		advanceForm.AddFormItem(formItem)
	}

	// add buttons
	var advanceFormButtons []*tview.Button
	advanceForm.AddButton("save", func() {
		v.newConfirmModal("do you want to save and restart?", func() {
			save(func() {
				v.restart()
			})
		}, nil)
	})
	advanceForm.AddButton("default", func() {
		v.newConfirmModal("do you want to reset to default and exit?", func() {
			err := cfg.ResetToDefault()
			if err != nil {
				v.CrossClipboard.ErrorChan <- xerror.NewRuntimeErrorf("can not reset config: %v", err)
				return
			}

			v.Stop()
			os.Exit(0)
		}, nil)
	})
	for i := 0; i < advanceForm.GetButtonCount(); i++ {
		advanceFormButtons = append(advanceFormButtons, advanceForm.GetButton(i))
	}
	advanceFormButtonsLen := len(advanceFormButtons)

	// set layout arrow input to change between setting
	settingLayout.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRight:
			v.app.SetFocus(advanceForm)
		case tcell.KeyLeft:
			v.app.SetFocus(settingForm)
		}
		return event
	})

	// set settingForm arrow up/down key
	settingForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyUp:
			idx, _ := settingForm.GetFocusedItemIndex()
			newIdx := idx - 1
			if newIdx < 0 || newIdx >= formItemsLen {
				return event
			}
			v.app.SetFocus(formItems[newIdx])
		case tcell.KeyDown:
			idx, _ := settingForm.GetFocusedItemIndex()
			newIdx := idx + 1
			if newIdx <= 0 || newIdx >= formItemsLen {
				return event
			}
			v.app.SetFocus(formItems[newIdx])
		}
		return event
	})

	// set settingForm arrow up/down key
	advanceForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyUp:
			idx, btnIdx := advanceForm.GetFocusedItemIndex()
			if btnIdx != -1 {
				newIdx := btnIdx - 1
				if newIdx < 0 {
					v.app.SetFocus(advanceFormItems[advanceFormItemsLen-1])
					break
				}
				v.app.SetFocus(advanceFormButtons[newIdx])
				break
			}
			newIdx := idx - 1
			if newIdx < 0 {
				break
			}
			v.app.SetFocus(advanceFormItems[newIdx])
		case tcell.KeyDown:
			idx, btnIdx := advanceForm.GetFocusedItemIndex()
			if btnIdx != -1 {
				newIdx := btnIdx + 1
				if newIdx >= advanceFormButtonsLen {
					break
				}
				v.app.SetFocus(advanceFormButtons[newIdx])
				break
			}
			newIdx := idx + 1
			if newIdx >= advanceFormItemsLen {
				v.app.SetFocus(advanceFormButtons[0])
				break
			}
			v.app.SetFocus(advanceFormItems[newIdx])
		}
		return event
	})

	return &Page{
		Title:   "Setting",
		Content: settingLayout,
	}
}

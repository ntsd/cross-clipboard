package ui

import (
	"github.com/rivo/tview"
)

const (
	yesLabel = "yes"
	noLabel  = "no"
)

func (v *View) newConfirmModal(label string, yesFunc func(), noFunc func()) {
	modal := tview.NewModal().
		SetText(label).
		AddButtons([]string{yesLabel, noLabel}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			switch buttonLabel {
			case yesLabel:
				if yesFunc != nil {
					yesFunc()
				}
			case noLabel:
				if noFunc != nil {
					noFunc()
				}
			}

			// set root back to layout after submitted
			v.app.SetRoot(v.layout, true).SetFocus(v.layout)
		})

	// set app root to modal
	v.app.SetRoot(modal, true).SetFocus(modal)
}

package ui

import (
	"fmt"

	"github.com/ntsd/cross-clipboard/pkg/device"
	"github.com/rivo/tview"
)

const (
	trustLabel = "Trust"
	blockLabel = "Block"
)

func (v *View) newTrustModal(dv *device.Device) {
	modal := tview.NewModal().
		SetText(fmt.Sprintf("%s (%s) want to connect, trust?", dv.Name, dv.OS)).
		AddButtons([]string{trustLabel, blockLabel}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			switch buttonLabel {
			case trustLabel:
				dv.Status = device.StatusConnected
			case blockLabel:
				dv.Status = device.StatusBlocked
			}

			v.CrossClipboard.DeviceManager.UpdateDevice(dv)

			// set root back to layout after submitted
			v.app.SetRoot(v.layout, true).SetFocus(v.layout)
		})

	// set app root to modal
	v.app.SetRoot(modal, true).SetFocus(modal)
}

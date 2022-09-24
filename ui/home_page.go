package ui

import (
	"github.com/rivo/tview"
)

func (v *View) newHomePage() *Page {
	flex := tview.NewGrid().
		AddItem(v.newClipboardBox(v.CrossClipboard, v.CrossClipboard.Config.HiddenText), 0, 0, 1, 1, 0, 0, false).
		AddItem(v.newDevicesBox(v.CrossClipboard), 0, 1, 1, 1, 0, 0, true)
	return &Page{
		Title:   "Home",
		Content: flex,
	}
}

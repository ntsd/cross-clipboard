package ui

import (
	"github.com/rivo/tview"
)

func (v *View) newHomePage() *Page {
	flex := tview.NewGrid().
		AddItem(v.newClipboardBox(), 0, 0, 1, 1, 0, 0, true).
		AddItem(v.newDevicesBox(), 0, 1, 1, 1, 0, 0, true)
	return &Page{
		Title:   "Home",
		Content: flex,
	}
}

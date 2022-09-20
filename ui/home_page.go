package ui

import (
	"github.com/rivo/tview"
)

func (v *View) NewHomePage() *Page {
	flex := tview.NewGrid().
		AddItem(v.clipboardBox(v.CrossClipboard, v.CrossClipboard.Config.HiddenText), 0, 0, 1, 1, 0, 0, true).
		AddItem(v.devicesBox(v.CrossClipboard), 0, 1, 1, 1, 0, 0, true)
	return &Page{
		Title:   "Home",
		Content: flex,
	}
}

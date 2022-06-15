package ui

import (
	"github.com/rivo/tview"
)

func ClipboardBox() tview.Primitive {
	return tview.NewBox().SetBorder(true).SetTitle("clipboards")
}

func DevicesBox() tview.Primitive {
	return tview.NewBox().SetBorder(true).SetTitle("devices")
}

func NewHomePage() *Page {
	flex := tview.NewFlex().
		AddItem(ClipboardBox(), 0, 2, true).
		AddItem(DevicesBox(), 30, 1, false)
	return &Page{
		Title:   "Home",
		Content: flex,
	}
}

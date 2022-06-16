package ui

import (
	"github.com/ntsd/cross-clipboard/pkg/cross_clipboard"
	"github.com/rivo/tview"
)

func ClipboardBox(cc *cross_clipboard.CrossClipboard) tview.Primitive {
	table := tview.NewTable().
		SetFixed(1, 1)

	go func() {
		for clipboards := range cc.Clipboard.ClipboardsChannel {
			for i, clipboard := range clipboards {
				table.SetCell(i, 0, tview.NewTableCell(string(clipboard)))
			}
		}
	}()

	table.SetBorder(true).SetTitle("clipboards")

	return table
}

func DevicesBox() tview.Primitive {
	return tview.NewBox().SetBorder(true).SetTitle("devices")
}

func (v *View) NewHomePage() *Page {
	flex := tview.NewFlex().
		AddItem(ClipboardBox(v.CrossClipboard), 0, 2, true).
		AddItem(DevicesBox(), 30, 1, false)
	return &Page{
		Title:   "Home",
		Content: flex,
	}
}

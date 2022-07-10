package ui

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/ntsd/cross-clipboard/pkg/crossclipboard"
	"github.com/rivo/tview"
)

func ClipboardBox(cc *crossclipboard.CrossClipboard) tview.Primitive {
	table := tview.NewTable().
		SetFixed(1, 1)

	go func() {
		for clipboards := range cc.ClipboardManager.ClipboardsChannel {
			table.Clear()
			table.SetCell(0, 0, tview.NewTableCell("text").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft))
			table.SetCell(0, 1, tview.NewTableCell("size").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
			table.SetCell(0, 2, tview.NewTableCell("time").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignRight))
			for i, clipboard := range clipboards {
				row := i + 1
				if clipboard.Size > 50 {
					table.SetCell(row, 0, tview.NewTableCell(string(clipboard.Data[:50])))
				} else {
					table.SetCell(row, 0, tview.NewTableCell(string(clipboard.Data)))
				}
				table.SetCell(row, 1, tview.NewTableCell(strconv.FormatUint(uint64(clipboard.Size), 10)))
				table.SetCell(row, 2, tview.NewTableCell(clipboard.Time.Format("15:04:05")))
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

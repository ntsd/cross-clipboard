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
			table.SetCell(0, 0, tview.NewTableCell("time").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft))
			table.SetCell(0, 1, tview.NewTableCell("bytes").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft))
			table.SetCell(0, 2, tview.NewTableCell("data").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft))

			for i, clipboard := range clipboards {
				row := i + 1
				table.SetCell(row, 0, tview.NewTableCell(clipboard.Time.Format("15:04:05")))
				table.SetCell(row, 1, tview.NewTableCell(strconv.FormatUint(uint64(clipboard.Size), 10)))
				if clipboard.Size > 50 {
					table.SetCell(row, 2, tview.NewTableCell(string(clipboard.Data[:50])))
				} else {
					table.SetCell(row, 2, tview.NewTableCell(string(clipboard.Data)))
				}
			}
		}
	}()

	table.SetBorder(true).SetTitle("clipboards")

	return table
}

func DevicesBox(cc *crossclipboard.CrossClipboard) tview.Primitive {
	table := tview.NewTable().
		SetFixed(1, 1)

	go func() {
		for devices := range cc.DeviceManager.DevicesChannel {
			table.Clear()
			table.SetCell(0, 0, tview.NewTableCell("id").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft))
			table.SetCell(0, 1, tview.NewTableCell("address").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft))

			row := 1
			for id, device := range devices {
				if len(id) > 10 {
					table.SetCell(row, 0, tview.NewTableCell(id[:10]))
				} else {
					table.SetCell(row, 0, tview.NewTableCell(id))
				}
				table.SetCell(row, 1, tview.NewTableCell(device.AddressInfo.Addrs[0].String()))
				row++
			}
		}
	}()

	table.SetBorder(true).SetTitle("devices")

	return table
}

func (v *View) NewHomePage() *Page {
	flex := tview.NewFlex().
		AddItem(ClipboardBox(v.CrossClipboard), 0, 2, true).
		AddItem(DevicesBox(v.CrossClipboard), 40, 1, false)
	return &Page{
		Title:   "Home",
		Content: flex,
	}
}

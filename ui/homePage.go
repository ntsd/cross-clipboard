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
				if clipboard.IsImage {
					table.SetCell(row, 2, tview.NewTableCell("image"))
					continue
				}
				table.SetCell(row, 2, tview.NewTableCell(limitTextLength(string(clipboard.Data), 10)))
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
			table.SetCell(0, 0, tview.NewTableCell("name").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft))
			table.SetCell(0, 1, tview.NewTableCell("address").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft))

			row := 1
			for id, device := range devices {
				name := device.Name
				if name == "" {
					name = id
				}
				table.SetCell(row, 0, tview.NewTableCell(limitTextLength(name, 10)))
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

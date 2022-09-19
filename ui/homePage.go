package ui

import (
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/ntsd/cross-clipboard/pkg/crossclipboard"
	"github.com/rivo/tview"
)

func ClipboardBox(cc *crossclipboard.CrossClipboard, hiddenText bool) tview.Primitive {
	table := tview.NewTable().
		SetFixed(1, 1)

	go func() {
		for clipboards := range cc.ClipboardManager.ClipboardsChannel {
			table.Clear()
			table.SetCell(0, 0, tview.NewTableCell("time").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft))
			table.SetCell(0, 1, tview.NewTableCell("bytes").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft))
			table.SetCell(0, 2, tview.NewTableCell("type").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft))
			if !hiddenText {
				table.SetCell(0, 3, tview.NewTableCell("text").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft))
			}

			for i, clipboard := range clipboards {
				row := i + 1
				table.SetCell(row, 0, tview.NewTableCell(clipboard.Time.Format("15:04:05")))
				table.SetCell(row, 1, tview.NewTableCell(strconv.FormatUint(uint64(clipboard.Size), 10)))
				if clipboard.IsImage {
					table.SetCell(row, 2, tview.NewTableCell("image"))
				} else {
					table.SetCell(row, 2, tview.NewTableCell("text"))
					if !hiddenText {
						table.SetCell(row, 3, tview.NewTableCell(limitTextLength(string(clipboard.Data), 10)))
					}
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
			table.SetCell(0, 0, tview.NewTableCell("name").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft))
			table.SetCell(0, 1, tview.NewTableCell("status").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft))
			table.SetCell(0, 2, tview.NewTableCell("address").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft))

			row := 1
			for id, device := range devices {
				name := device.Name
				if name == "" {
					name = id
				}
				table.SetCell(row, 0, tview.NewTableCell(limitTextLength(name, 10)))
				table.SetCell(row, 1, tview.NewTableCell(device.Status.ToString()))

				addressStr := ""
				for _, address := range device.AddressInfo.Addrs {
					if !strings.Contains(address.String(), "127.0.0.1") {
						addressStr = address.String()
					}
				}

				table.SetCell(row, 2, tview.NewTableCell(addressStr))
				row++
			}
		}
	}()

	table.SetBorder(true).SetTitle("devices")

	return table
}

func (v *View) NewHomePage() *Page {
	flex := tview.NewGrid().
		AddItem(ClipboardBox(v.CrossClipboard, v.CrossClipboard.Config.HiddenText), 0, 0, 1, 1, 0, 0, true).
		AddItem(DevicesBox(v.CrossClipboard), 0, 1, 1, 1, 0, 0, true)
	return &Page{
		Title:   "Home",
		Content: flex,
	}
}

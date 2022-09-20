package ui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/ntsd/cross-clipboard/pkg/crossclipboard"
	"github.com/ntsd/cross-clipboard/pkg/device"
	"github.com/rivo/tview"
)

func (v *View) devicesBox(cc *crossclipboard.CrossClipboard) tview.Primitive {
	table := tview.NewTable().
		SetFixed(1, 1)

	go func() {
		for devices := range cc.DeviceManager.DevicesChannel {
			table.Clear()
			table.SetCell(0, 0, tview.NewTableCell("name").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft))
			table.SetCell(0, 1, tview.NewTableCell("status").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft))
			table.SetCell(0, 2, tview.NewTableCell("address").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft))

			row := 1
			for id, dv := range devices {
				if dv.Status == device.StatusPending {
					v.newTrustModal(dv)
				}

				name := dv.Name
				if name == "" {
					name = id
				}
				table.SetCell(row, 0, tview.NewTableCell(limitTextLength(name, 10)))
				table.SetCell(row, 1, tview.NewTableCell(dv.Status.ToString()))

				addressStr := ""
				for _, address := range dv.AddressInfo.Addrs {
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

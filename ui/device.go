package ui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/ntsd/cross-clipboard/pkg/device"
	"github.com/ntsd/go-utils/pkg/stringutil"
	"github.com/rivo/tview"
)

var deviceStatusTextColorMap = map[device.DeviceStatus]tcell.Color{
	device.StatusPending:      tcell.ColorYellow,
	device.StatusConnected:    tcell.ColorGreen,
	device.StatusDisconnected: tcell.ColorGray,
	device.StatusError:        tcell.ColorRed,
	device.StatusBlocked:      tcell.ColorRed,
}

func (v *View) newDevicesBox() tview.Primitive {
	table := tview.NewTable().
		SetFixed(1, 1)

	cc := v.CrossClipboard

	go func() {
		for devices := range cc.DeviceManager.DevicesChannel {
			table.Clear()
			v.renderDevicesTable(table, devices)
		}
	}()

	table.SetBorder(true).SetTitle("devices")

	return table
}

func (v *View) renderDevicesTable(table *tview.Table, devices map[string]*device.Device) {
	headerColor := tcell.ColorYellow
	table.SetCell(0, 0, tview.NewTableCell("name").SetTextColor(headerColor).SetAlign(tview.AlignLeft))
	table.SetCell(0, 1, tview.NewTableCell("status").SetTextColor(headerColor).SetAlign(tview.AlignLeft))
	table.SetCell(0, 2, tview.NewTableCell("os").SetTextColor(headerColor).SetAlign(tview.AlignLeft))
	table.SetCell(0, 3, tview.NewTableCell("address").SetTextColor(headerColor).SetAlign(tview.AlignLeft))

	row := 1
	for id, dv := range devices {
		if dv.Status == device.StatusPending {
			v.newTrustModal(dv)
		}

		textColor, ok := deviceStatusTextColorMap[dv.Status]
		if !ok {
			textColor = tcell.ColorDefault
		}

		name := dv.Name
		if name == "" {
			name = id
		}
		table.SetCell(row, 0, tview.NewTableCell(stringutil.LimitStringLen(name, 10)))
		table.SetCell(row, 1, tview.NewTableCell(string(dv.Status)).SetTextColor(textColor))
		table.SetCell(row, 2, tview.NewTableCell(dv.OS))

		// because multiaddr sometimes include 127.0.0.1
		addressStr := ""
		for _, address := range dv.AddressInfo.Addrs {
			if !strings.Contains(address.String(), "127.0.0.1") {
				addressStr = address.String()
			}
		}

		table.SetCell(row, 3, tview.NewTableCell(addressStr))
		row++
	}
}

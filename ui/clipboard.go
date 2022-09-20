package ui

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/ntsd/cross-clipboard/pkg/crossclipboard"
	"github.com/rivo/tview"
)

func (v *View) newClipboardBox(cc *crossclipboard.CrossClipboard, hiddenText bool) tview.Primitive {
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

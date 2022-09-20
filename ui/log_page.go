package ui

import (
	"fmt"

	"github.com/rivo/tview"
)

func (v *View) NewLogPage() *Page {
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			v.app.Draw()
		})

	go func() {
		for {
			select {
			case log := <-v.CrossClipboard.LogChan:
				fmt.Fprint(textView, fmt.Sprintf("[blue]log:[white] %s\n", log))

			case err := <-v.CrossClipboard.ErrChan:
				fmt.Fprint(textView, fmt.Sprintf("[red]err: %s[white]\n", err))
			}
		}
	}()

	return &Page{
		Title:   "Log",
		Content: textView,
	}
}

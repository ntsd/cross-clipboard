package ui

import (
	"errors"
	"fmt"
	"log"

	"github.com/ntsd/cross-clipboard/pkg/xerror"
	"github.com/rivo/tview"
)

func (v *View) newLogPage() *Page {
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
			case err := <-v.CrossClipboard.ErrorChan:
				var fatalErr *xerror.FatalError
				if errors.As(err, &fatalErr) {
					v.app.Stop()
					log.Fatal(fmt.Errorf("fatal error: %w", fatalErr))
				}

				fmt.Fprint(textView, fmt.Sprintf("[red]err: %s[white]\n", err))
			}
		}
	}()

	return &Page{
		Title:   "Log",
		Content: textView,
	}
}

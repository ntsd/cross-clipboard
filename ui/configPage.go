package ui

import (
	"fmt"

	"github.com/rivo/tview"
)

func (v *View) NewConfigPage() *Page {
	textView := tview.NewTextView()
	fmt.Fprint(textView, "ID: "+v.CrossClipboard.Host.ID().Pretty())
	return &Page{
		Title:   "Config",
		Content: textView,
	}
}

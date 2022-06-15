package ui

import (
	"fmt"

	"github.com/rivo/tview"
)

func NewConfigPage() *Page {
	textView := tview.NewTextView()
	fmt.Fprint(textView, "Config Page")
	return &Page{
		Title:   "Config",
		Content: textView,
	}
}

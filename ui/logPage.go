package ui

import (
	"fmt"

	"github.com/rivo/tview"
)

func NewLogPage() *Page {
	textView := tview.NewTextView()
	fmt.Fprint(textView, "Log Page")
	return &Page{
		Title:   "Log",
		Content: textView,
	}
}

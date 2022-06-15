package ui

import (
	"fmt"

	"github.com/rivo/tview"
)

func ConfigPage(prevPage func(), nextPage func()) (title string, content tview.Primitive) {
	textView := tview.NewTextView()
	fmt.Fprint(textView, "Config Page")
	return "Config", textView
}

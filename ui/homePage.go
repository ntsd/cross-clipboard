package ui

import (
	"fmt"

	"github.com/rivo/tview"
)

func HomePage(prevPage func(), nextPage func()) (title string, content tview.Primitive) {
	textView := tview.NewTextView()
	fmt.Fprint(textView, "Home Page")
	return "Home", textView
}

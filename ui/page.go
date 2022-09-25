package ui

import (
	"strconv"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Page struct {
	Title   string
	Content tview.Primitive
}

func (v *View) pageInputCapture(event *tcell.EventKey) *tcell.EventKey {
	// TODO prevent form type to not change the page

	if unicode.IsDigit(event.Rune()) {
		pageNum := int(event.Rune() - '1')
		if pageNum < len(v.pages) && pageNum >= 0 {
			v.goToPage(pageNum)
		}
	}
	return event
}

func (v *View) goToPage(pageNum int) {
	v.menuBar.Highlight(strconv.Itoa(pageNum)).
		ScrollToHighlight()
}

func (v *View) previousPage() {
	page, _ := strconv.Atoi(v.menuBar.GetHighlights()[0])
	page = (page - 1 + len(v.pages)) % len(v.pages)
	v.menuBar.Highlight(strconv.Itoa(page)).
		ScrollToHighlight()
}

func (v *View) nextPage() {
	currentPage, _ := strconv.Atoi(v.menuBar.GetHighlights()[0])
	newPage := (currentPage + 1) % len(v.pages)
	v.menuBar.Highlight(strconv.Itoa(newPage)).
		ScrollToHighlight()
}

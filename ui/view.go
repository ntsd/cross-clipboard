package ui

import (
	"fmt"
	"strconv"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/ntsd/cross-clipboard/pkg/cross_clipboard"
	"github.com/rivo/tview"
)

type View struct {
	CrossClipboard *cross_clipboard.CrossClipboard
	app            *tview.Application
	layout         *tview.Flex

	info      *tview.TextView
	basePages *tview.Pages
	pages     []Page
}

func NewView(cc *cross_clipboard.CrossClipboard) *View {
	view := &View{
		CrossClipboard: cc,
		layout:         tview.NewFlex(),
		info:           tview.NewTextView(),
		basePages:      tview.NewPages(),
		pages: []Page{
			HomePage,
			ConfigPage,
		},
	}

	return view
}

func (v *View) Start() {
	v.app = tview.NewApplication()

	// Set bottom info bar
	v.info.SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false).
		SetHighlightedFunc(func(added, removed, remaining []string) {
			v.basePages.SwitchToPage(added[0])
		})

	// Set layout with basePages and info
	v.layout.SetDirection(tview.FlexRow).
		AddItem(v.basePages, 0, 1, true).
		AddItem(v.info, 1, 1, false)

	// Create pages and controller
	goToPage := func(pageNum int) {
		v.info.Highlight(strconv.Itoa(pageNum)).
			ScrollToHighlight()
	}
	previousPage := func() {
		page, _ := strconv.Atoi(v.info.GetHighlights()[0])
		page = (page - 1 + len(v.pages)) % len(v.pages)
		v.info.Highlight(strconv.Itoa(page)).
			ScrollToHighlight()
	}
	nextPage := func() {
		currentPage, _ := strconv.Atoi(v.info.GetHighlights()[0])
		newPage := (currentPage + 1) % len(v.pages)
		v.info.Highlight(strconv.Itoa(newPage)).
			ScrollToHighlight()
	}
	for index, page := range v.pages {
		title, primitive := page(previousPage, nextPage)
		v.basePages.AddPage(strconv.Itoa(index), primitive, true, index == 0)
		fmt.Fprintf(v.info, `%d ["%d"][darkcyan]%s[white][""]  `, index+1, index, title)
	}
	v.info.Highlight("0")

	v.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if unicode.IsDigit(event.Rune()) {
			goToPage(int(event.Rune() - '1'))
		}
		return event
	})

	if err := v.app.
		SetRoot(v.layout, true).
		EnableMouse(true).
		Run(); err != nil {
		panic(err)
	}
}

func (v *View) Stop() {
	v.app.Stop()
}

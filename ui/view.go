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
	pages     []*Page
}

func NewView(cc *cross_clipboard.CrossClipboard) *View {
	v := &View{
		CrossClipboard: cc,
		layout:         tview.NewFlex(),
		info:           tview.NewTextView(),
		basePages:      tview.NewPages(),
	}

	v.pages = []*Page{
		v.NewHomePage(),
		v.NewConfigPage(),
		v.NewLogPage(),
	}

	return v
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

	// Create pages
	for index, page := range v.pages {
		v.basePages.AddPage(strconv.Itoa(index), page.Content, true, index == 0)
		fmt.Fprintf(v.info, `%d ["%d"][darkcyan]%s[white][""]  `, index+1, index, page.Title)
	}
	v.info.Highlight("0")

	v.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if unicode.IsDigit(event.Rune()) {
			v.goToPage(int(event.Rune() - '1'))
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

func (v *View) goToPage(pageNum int) {
	v.info.Highlight(strconv.Itoa(pageNum)).
		ScrollToHighlight()
}

func (v *View) previousPage() {
	page, _ := strconv.Atoi(v.info.GetHighlights()[0])
	page = (page - 1 + len(v.pages)) % len(v.pages)
	v.info.Highlight(strconv.Itoa(page)).
		ScrollToHighlight()
}

func (v *View) nextPage() {
	currentPage, _ := strconv.Atoi(v.info.GetHighlights()[0])
	newPage := (currentPage + 1) % len(v.pages)
	v.info.Highlight(strconv.Itoa(newPage)).
		ScrollToHighlight()
}

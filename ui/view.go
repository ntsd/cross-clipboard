package ui

import (
	"fmt"
	"strconv"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/ntsd/cross-clipboard/pkg/crossclipboard"
	"github.com/rivo/tview"
)

type View struct {
	CrossClipboard *crossclipboard.CrossClipboard
	app            *tview.Application
	layout         *tview.Flex

	menuBar   *tview.TextView
	basePages *tview.Pages
	pages     []*Page
}

func NewView(cc *crossclipboard.CrossClipboard) *View {
	v := &View{
		CrossClipboard: cc,
		layout:         tview.NewFlex(),
		menuBar:        tview.NewTextView(),
		basePages:      tview.NewPages(),
	}

	v.pages = []*Page{
		v.newHomePage(),
		v.newSettingPage(),
		v.newLogPage(),
	}

	return v
}

func (v *View) Start() {
	v.app = tview.NewApplication()

	// Set bottom menu bar
	v.menuBar.SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false).
		SetHighlightedFunc(func(added, removed, remaining []string) {
			v.basePages.SwitchToPage(added[0])
		})

	// Set layout with basePages and info
	v.layout.SetDirection(tview.FlexRow).
		AddItem(v.basePages, 0, 1, true).
		AddItem(v.menuBar, 1, 1, false)

	// Create pages
	for index, page := range v.pages {
		v.basePages.AddPage(strconv.Itoa(index), page.Content, true, index == 0)
		fmt.Fprintf(v.menuBar, `%d ["%d"][darkcyan]%s[white][""]  `, index+1, index, page.Title)
	}
	v.menuBar.Highlight("0")

	v.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// check form text input to avoid menu changing
		focusFormIdx, _ := v.pages[1].Content.(*tview.Form).GetFocusedItemIndex()
		if contains([]int{1, 2, 5, 6, 7}, focusFormIdx) {
			return event
		}

		if unicode.IsDigit(event.Rune()) {
			pageNum := int(event.Rune() - '1')
			if pageNum < len(v.pages) {
				v.goToPage(pageNum)
			}
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

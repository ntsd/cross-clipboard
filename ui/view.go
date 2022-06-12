package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type View struct {
	app    *tview.Application
	layout *tview.Grid

	basePages *tview.Pages
	mainPage  tview.Primitive
}

func NewView() *View {
	view := &View{
		layout:    newGrid(),
		basePages: tview.NewPages(),
		mainPage:  mainPageGrid(),
	}

	return view
}

func (v *View) Start() {
	v.app = tview.NewApplication()

	v.app.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		screen.Clear()
		return false
	})

	v.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return event
	})

	v.basePages.AddPage("main-page", v.mainPage, true, true)
	v.layout.AddItem(v.basePages, 1, 0, 1, 1, 0, 0, true)

	if err := v.app.
		SetRoot(v.layout, true).
		SetFocus(v.layout).
		Run(); err != nil {
		panic(err)
	}
}

func (v *View) Stop() {
	v.app.Stop()
}

func newGrid() *tview.Grid {
	return tview.NewGrid().
		SetRows(1, 0).
		SetColumns(0)
}

func mainPageGrid() *tview.Grid {
	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}

	menu := newPrimitive("Menu")
	main := newPrimitive("Main content")

	grid := tview.NewGrid().
		SetRows(2, 1).
		SetBorders(true)

	// Layout for screens narrower than 100 cells (menu and side bar are hidden).
	grid.AddItem(menu, 0, 0, 0, 0, 0, 0, false).
		AddItem(main, 1, 0, 1, 3, 0, 0, false)

	// Layout for screens wider than 100 cells.
	grid.AddItem(menu, 1, 0, 1, 1, 0, 100, false).
		AddItem(main, 1, 1, 1, 1, 0, 100, false)

	return grid
}

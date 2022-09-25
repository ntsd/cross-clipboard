package ui

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/ntsd/cross-clipboard/pkg/config"
	"github.com/ntsd/cross-clipboard/pkg/crossclipboard"
	"github.com/ntsd/cross-clipboard/pkg/xerror"
	"github.com/rivo/tview"
)

type View struct {
	CrossClipboard *crossclipboard.CrossClipboard
	app            *tview.Application

	layout    *tview.Flex
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

	// Add pages and menu bar
	for index, page := range v.pages {
		v.basePages.AddPage(strconv.Itoa(index), page.Content, true, index == 0)
		fmt.Fprintf(v.menuBar, `%d ["%d"][darkcyan]%s[white][""]  `, index+1, index, page.Title)
	}
	v.menuBar.Highlight("0")

	// handle shortcuts key
	v.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		return v.pageInputCapture(event)
	})

	if err := v.app.
		SetRoot(v.layout, true).
		EnableMouse(true).
		Run(); err != nil {
		panic(err)
	}
}

func (v *View) Stop() {
	v.CrossClipboard.Stop()
	v.app.Stop()
	v = nil
}

func (v *View) restart() {
	v.Stop()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	crossClipboard, err := crossclipboard.NewCrossClipboard(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.TerminalMode {
		for {
			select {
			case l := <-crossClipboard.LogChan:
				log.Println("log: ", l)
			case err := <-crossClipboard.ErrorChan:
				var fatalErr *xerror.FatalError
				if errors.As(err, &fatalErr) {
					log.Fatal(fmt.Errorf("fatal error: %w", fatalErr))
				}
				log.Println(fmt.Errorf("runtime error: %w", err))
			case cb := <-crossClipboard.ClipboardManager.ClipboardsChannel:
				_ = cb
			case dv := <-crossClipboard.DeviceManager.DevicesChannel:
				_ = dv
			}
		}
	} else {
		v = NewView(crossClipboard)
		v.Start()
	}
}

package main

import (
	"github.com/ntsd/cross-clipboard/pkg/config"
	"github.com/ntsd/cross-clipboard/pkg/p2p"
	"github.com/ntsd/cross-clipboard/ui/assets"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	APP_ID      = "dev.ntsd.cross.clipboard"
	WINDOW_NAME = "Cross Clipbaord"
)

func main() {
	cfg := config.Config{
		RendezvousString: "default-group",
		ProtocolID:       "/cross-clipboard/0.0.1",
		ListenHost:       "0.0.0.0",
		ListenPort:       4001,
	}
	p2p.StartP2P(cfg)

	a := app.NewWithID(APP_ID)
	a.SetIcon(assets.ResourceLogoPng)

	w := a.NewWindow(WINDOW_NAME)

	hello := widget.NewLabel("Hello Fyne!")
	vBox := container.NewVBox(
		hello,
		widget.NewButton("Hi!", func() {
			hello.SetText("Welcome :)")
		}),
	)
	vScroll := container.NewVScroll(vBox)

	w.SetContent(vScroll)

	w.ShowAndRun()
}

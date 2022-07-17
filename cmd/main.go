package main

import (
	"log"

	"github.com/ntsd/cross-clipboard/pkg/config"
	"github.com/ntsd/cross-clipboard/pkg/crossclipboard"
	"github.com/ntsd/cross-clipboard/ui"
)

func main() {
	cfg := config.LoadConfig()

	crossClipboard, err := crossclipboard.NewCrossClipboard(cfg)
	if err != nil {
		panic(err)
	}

	if cfg.TerminalMode {
		for {
			select {
			case l := <-crossClipboard.LogChan:
				log.Println("log: ", l)
			case err := <-crossClipboard.ErrChan:
				log.Panicln("err: ", err)
			case cb := <-crossClipboard.ClipboardManager.ClipboardsChannel:
				_ = cb
			case dv := <-crossClipboard.DeviceManager.DevicesChannel:
				_ = dv
			}
		}
	} else {
		view := ui.NewView(crossClipboard)
		view.Start()
	}
}

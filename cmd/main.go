package main

import (
	"log"

	"github.com/ntsd/cross-clipboard/pkg/config"
	"github.com/ntsd/cross-clipboard/pkg/cross_clipboard"
	"github.com/ntsd/cross-clipboard/ui"
)

func main() {
	cfg := config.LoadConfig()

	crossClipboard, err := cross_clipboard.NewCrossClipboard(cfg)
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
			}
		}
	} else {
		view := ui.NewView(crossClipboard)
		view.Start()
	}
}

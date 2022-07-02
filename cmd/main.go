package main

import (
	"log"

	"github.com/ntsd/cross-clipboard/pkg/config"
	"github.com/ntsd/cross-clipboard/pkg/cross_clipboard"
	"github.com/ntsd/cross-clipboard/ui"
)

func main() {
	// help := flag.Bool("help", false, "Display Help")
	// if *help {
	// 	fmt.Printf("Simple example for peer discovery using mDNS. mDNS is great when you have multiple peers in local LAN.")
	// 	fmt.Printf("Usage: \n   Run './chat-with-mdns'\nor Run './chat-with-mdns -host [host] -port [port] -rendezvous [string] -pid [proto ID]'\n")

	// 	os.Exit(0)
	// }

	// cfg := utils.ParseFlags()

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

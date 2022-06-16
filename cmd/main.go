package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ntsd/cross-clipboard/pkg/cross_clipboard"
	"github.com/ntsd/cross-clipboard/pkg/utils"
)

func main() {
	help := flag.Bool("help", false, "Display Help")
	if *help {
		fmt.Printf("Simple example for peer discovery using mDNS. mDNS is great when you have multiple peers in local LAN.")
		fmt.Printf("Usage: \n   Run './chat-with-mdns'\nor Run './chat-with-mdns -host [host] -port [port] -rendezvous [string] -pid [proto ID]'\n")

		os.Exit(0)
	}

	cfg := utils.ParseFlags()

	crossClipboard, err := cross_clipboard.NewCrossClipboard(cfg)
	if err != nil {
		panic(err)
	}
	_ = crossClipboard

	// view := ui.NewView(crossClipboard)
	// view.Start()

	for {
		select {
		case log := <-crossClipboard.LogChan:
			fmt.Println("log: ", log)

		case err := <-crossClipboard.ErrChan:
			fmt.Println("err: ", err)
		}
	}
}

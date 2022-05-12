package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ntsd/cross-clipboard/pkg/p2p"
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
	p2p.StartP2P(cfg)
}

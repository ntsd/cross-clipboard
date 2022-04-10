package utils

import (
	"flag"
)

type config struct {
	RendezvousString string
	ProtocolID       string
	ListenHost       string
	ListenPort       int
}

func ParseFlags() *config {
	c := &config{}

	flag.StringVar(&c.RendezvousString, "rendezvous", "meetme", "Unique string to identify group of nodes. Share this with your friends to let them connect with you")
	flag.StringVar(&c.ListenHost, "host", "0.0.0.0", "The bootstrap node host listen address\n")
	flag.StringVar(&c.ProtocolID, "pid", "/chat/1.1.0", "Sets a protocol id for stream headers")
	flag.IntVar(&c.ListenPort, "port", 4001, "node listen port")

	flag.Parse()
	return c
}

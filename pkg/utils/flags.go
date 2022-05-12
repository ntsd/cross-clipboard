package utils

import (
	"flag"
)

type Config struct {
	RendezvousString string
	ProtocolID       string
	ListenHost       string
	ListenPort       int
}

func ParseFlags() Config {
	c := Config{}

	flag.StringVar(&c.RendezvousString, "rendezvous", "default-group", "Unique string to identify group of nodes. Share this with your friends to let them connect with you")
	flag.StringVar(&c.ProtocolID, "pid", "/cross-clipboard/0.0.1", "Sets a protocol id for stream headers")
	flag.StringVar(&c.ListenHost, "host", "0.0.0.0", "The bootstrap node host listen address\n")
	flag.IntVar(&c.ListenPort, "port", 4001, "node listen port")

	flag.Parse()
	return c
}

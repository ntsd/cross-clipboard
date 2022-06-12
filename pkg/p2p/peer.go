package p2p

import (
	"bufio"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
)

type Peer struct {
	AddressInfo peer.AddrInfo
	OS          string

	Stream network.Stream
	Writer *bufio.Writer
	Reader *bufio.Reader

	Status    string
	PublicKey string

	LogChan chan string
	ErrChan chan error
}

func NewPeer(
	addrInfo peer.AddrInfo,
	stream network.Stream,
) *Peer {
	return &Peer{
		AddressInfo: addrInfo,
		Stream:      stream,
		Reader:      bufio.NewReader(stream),
		Writer:      bufio.NewWriter(stream),
	}
}

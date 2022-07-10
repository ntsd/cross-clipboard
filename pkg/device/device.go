package device

import (
	"bufio"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
)

// Device struct for peer
type Device struct {
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

// NewDevice initial new peer
func NewDevice(
	addrInfo peer.AddrInfo,
	stream network.Stream,
) *Device {
	return &Device{
		AddressInfo: addrInfo,
		Stream:      stream,
		Reader:      bufio.NewReader(stream),
		Writer:      bufio.NewWriter(stream),
	}
}

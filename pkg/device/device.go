package device

import (
	"bufio"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
)

// DeviceStatus device status
type DeviceStatus int

const (
	// StatusPending the device waiting to handshake and trust the device
	StatusPending DeviceStatus = iota
	// StatusConnecting the device is trusted and connecting
	StatusConnecting
	// StatusConnecting the device is trusted but disconnecting or offline
	StatusDisconnecting
	// StatusError found a error in the device should disconnect and reconnect
	StatusError
	// StatusBlocked the device is blocked by the user
	StatusBlocked
)

// Device struct for peer
type Device struct {
	AddressInfo peer.AddrInfo

	OS        string
	Name      string
	PublicKey *[]byte
	Status    DeviceStatus

	Stream network.Stream
	Writer *bufio.Writer
	Reader *bufio.Reader

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

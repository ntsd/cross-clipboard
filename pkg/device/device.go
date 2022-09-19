package device

import (
	"bufio"
	"fmt"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/ntsd/cross-clipboard/pkg/crypto"
	"github.com/ntsd/cross-clipboard/pkg/protobuf"
)

// Device struct for peer
type Device struct {
	AddressInfo peer.AddrInfo

	OS        string
	Name      string
	PublicKey []byte
	Status    DeviceStatus

	Stream network.Stream
	Writer *bufio.Writer
	Reader *bufio.Reader

	LogChan chan string
	ErrChan chan error

	PgpEncrypter *crypto.PGPEncrypter
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

// UpdateFromProtobuf update device from protobuf device data
func (dv *Device) UpdateFromProtobuf(deviceData *protobuf.DeviceData) error {
	dv.Name = deviceData.Name
	dv.OS = deviceData.Os
	dv.PublicKey = deviceData.PublicKey

	publicKey, err := crypto.ByteToPGPKey(deviceData.PublicKey)
	if err != nil {
		return fmt.Errorf("error to create pgp public key: %w", err)
	}
	pgpEncrypter, err := crypto.NewPGPEncrypter(publicKey)
	if err != nil {
		return fmt.Errorf("error to create pgp encrypter: %w", err)
	}
	dv.PgpEncrypter = pgpEncrypter

	return nil
}

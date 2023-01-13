package device

import (
	"bufio"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/ntsd/cross-clipboard/pkg/crypto"
	"github.com/ntsd/cross-clipboard/pkg/protobuf"
	"github.com/ntsd/cross-clipboard/pkg/xerror"
)

// Device struct for peer
type Device struct {
	AddressInfo peer.AddrInfo `json:"-"`

	OS        string       `json:"os"`
	Name      string       `json:"name"`
	PublicKey []byte       `json:"publicKey"`
	Status    DeviceStatus `json:"status"`

	Stream network.Stream `json:"-"`
	Writer *bufio.Writer  `json:"-"`
	Reader *bufio.Reader  `json:"-"`

	PgpEncrypter *crypto.PGPEncrypter `json:"-"`
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

// Trust trust this device and change status to connected
func (dv *Device) Trust() error {
	err := dv.CreatePGPEncrypter()
	if err != nil {
		return xerror.NewRuntimeError("can not create pgp encrypter").Wrap(err)
	}

	dv.Status = StatusConnected
	return nil
}

// Block block this device
func (dv *Device) Block() {
	dv.Status = StatusBlocked
}

// UpdateFromProtobuf update device from protobuf device data
func (dv *Device) UpdateFromProtobuf(deviceData *protobuf.DeviceData) {
	dv.Name = deviceData.Name
	dv.OS = deviceData.Os
	dv.PublicKey = deviceData.PublicKey
}

func (dv *Device) CreatePGPEncrypter() error {
	publicKey, err := crypto.ByteToPGPKey(dv.PublicKey)
	if err != nil {
		return xerror.NewRuntimeError("error to create pgp public key").Wrap(err)
	}
	pgpEncrypter, err := crypto.NewPGPEncrypter(publicKey)
	if err != nil {
		return xerror.NewRuntimeError("error to create pgp encrypter").Wrap(err)
	}
	dv.PgpEncrypter = pgpEncrypter

	return nil
}

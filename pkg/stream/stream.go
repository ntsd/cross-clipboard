package stream

import (
	"bufio"
	"fmt"
	"runtime"
	"time"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/ntsd/cross-clipboard/pkg/clipboard"
	"github.com/ntsd/cross-clipboard/pkg/config"
	"github.com/ntsd/cross-clipboard/pkg/crypto"
	"github.com/ntsd/cross-clipboard/pkg/device"
	"github.com/ntsd/cross-clipboard/pkg/devicemanager"
)

const (
	// EOF is the end of message byte use to delim the message
	EOF byte = 0x00
	// DATA TYPE is the last byte befor EOF use to determine the message type
	DATA_TYPE_DEVICE    byte = 0xFF
	DATA_TYPE_CLIPBOARD byte = 0xFE
)

// StreamHandler struct for stream handler
type StreamHandler struct {
	Config           config.Config
	ClipboardManager *clipboard.ClipboardManager
	DeviceManager    *devicemanager.DeviceManager
	LogChan          chan string
	ErrorChan        chan error

	pgpDecrypter *crypto.PGPDecrypter
}

// NewStreamHandler initial new stream handler
func NewStreamHandler(
	cfg config.Config,
	cp *clipboard.ClipboardManager,
	deviceManager *devicemanager.DeviceManager,
	logChan chan string,
	errorChan chan error,
	pgpDecrypter *crypto.PGPDecrypter,
) *StreamHandler {
	s := &StreamHandler{
		Config:           cfg,
		ClipboardManager: cp,
		DeviceManager:    deviceManager,
		LogChan:          logChan,
		ErrorChan:        errorChan,
		pgpDecrypter:     pgpDecrypter,
	}
	go s.CreateWriteData()
	return s
}

// HandleStream handler when a peer connect this host
func (s *StreamHandler) HandleStream(stream network.Stream) {
	s.LogChan <- fmt.Sprintf("peer %s connecting to this host", stream.Conn().RemotePeer())

	// Create a new peer
	dv := device.NewDevice(peer.AddrInfo{
		ID:    stream.Conn().RemotePeer(),
		Addrs: []multiaddr.Multiaddr{stream.Conn().RemoteMultiaddr()},
	}, stream)
	s.DeviceManager.AddDevice(dv)
	dv.Reader = bufio.NewReader(stream)
	dv.Writer = bufio.NewWriter(stream)
	go s.CreateReadData(dv.Reader, stream.Conn().RemotePeer().Pretty())

	s.LogChan <- fmt.Sprintf("peer %s connected to this host", stream.Conn().RemotePeer())
	// 'stream' will stay open until you close it (or the other side closes it).
}

// CreateReadData craete a new read streaming for host or peer
func (s *StreamHandler) CreateReadData(reader *bufio.Reader, id string) {
	// generate public key
	pub, err := s.Config.PGPPrivateKey.GetPublicKey()
	if err != nil {
		s.ErrorChan <- fmt.Errorf("error to generate pubic key: %w", err)
		return
	}

	// sending device info and public key
	dv := s.DeviceManager.GetDevice(id)
	s.EncodeDeviceData(&DeviceData{
		Name:      s.Config.Username,
		Os:        runtime.GOOS,
		PublicKey: pub,
	})

	for {
		bytes, err := reader.ReadBytes(EOF)
		if err != nil {
			s.ErrorChan <- fmt.Errorf("error reading from buffer: %w", err)
			break
		}

		clipboardData, deviceData, err := s.DecodeData(bytes)
		if err != nil {
			s.ErrorChan <- fmt.Errorf("error decoding data: %w", err)
			continue
		}

		if clipboardData != nil {
			s.LogChan <- fmt.Sprintf("received clipboard data, peer: %s size: %d", id, clipboardData.DataSize)
			s.ClipboardManager.WriteClipboard(clipboard.Clipboard{
				IsImage: clipboardData.IsImage,
				Data:    clipboardData.Data,
				Size:    clipboardData.DataSize,
				Time:    time.UnixMicro(clipboardData.Time),
			})
		}

		if deviceData != nil {
			s.LogChan <- fmt.Sprintf("received device data, peer: %s", id)
			s.LogChan <- fmt.Sprintf("%s wanted to connect", deviceData.Name)

			dv.Name = deviceData.Name
			dv.OS = deviceData.Os
			dv.PublicKey = deviceData.PublicKey

			publicKey, err := crypto.ByteToPGPKey(deviceData.PublicKey)
			if err != nil {
				s.ErrorChan <- fmt.Errorf("error to create pgp public key: %w", err)
				break
			}
			pgpEncrypter, err := crypto.NewPGPEncrypter(publicKey)
			if err != nil {
				s.ErrorChan <- fmt.Errorf("error to create pgp encrypter: %w", err)
				break
			}
			dv.PgpEncrypter = pgpEncrypter

			dv.Status = device.StatusConnecting

			s.DeviceManager.UpdateDevice(dv)
		}
	}
	s.LogChan <- fmt.Sprintf("ending read stream for peer: %s", id)
}

// CreateWriteData handle clipboad channel and write to all peers and host
func (s *StreamHandler) CreateWriteData() {
	// waiting for clipboard data
	for clipboardBytes := range s.ClipboardManager.ReadChannel {
		length := len(clipboardBytes)
		if length == 0 {
			// ignore empty clipboard
			continue
		}

		now := time.Now()

		// set current clipbaord to avoid recursive
		s.ClipboardManager.AddClipboard(clipboard.Clipboard{
			Data: clipboardBytes,
			Size: uint32(length),
			Time: now,
		})

		clipboardData := &ClipboardData{
			IsImage:  false,
			Data:     clipboardBytes,
			DataSize: uint32(length),
			Time:     now.Unix(),
		}

		// send data to each devices
		for name, d := range s.DeviceManager.Devices {
			if d.PgpEncrypter == nil {
				s.ErrorChan <- fmt.Errorf("not found pgp encrypter for device %s", name)
				continue
			}

			s.LogChan <- fmt.Sprintf("sending data to peer: %s size: %d data: %s", name, length, string(clipboardBytes))

			clipboardDataBytes, err := s.EncodeClipboardData(name, clipboardData)
			if err != nil {
				s.ErrorChan <- fmt.Errorf("error encoding data: %w", err)
				continue
			}

			err = s.WriteData(d.Writer, clipboardDataBytes)
			if err != nil {
				s.LogChan <- fmt.Sprintf("ending write stream %s", name)
				s.DeviceManager.RemoveDevice(d)
			}
		}
	}
	s.LogChan <- fmt.Sprintf("ending write streams")
}

// WriteData write data to the writer
func (s *StreamHandler) WriteData(w *bufio.Writer, data []byte) error {
	_, err := w.Write(data)
	if err != nil {
		s.ErrorChan <- fmt.Errorf("error writing to buffer: %w", err)
		return err
	}

	err = w.Flush()
	if err != nil {
		s.ErrorChan <- fmt.Errorf("error flushing buffer: %w", err)
		return err
	}
	return nil
}

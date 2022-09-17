package stream

import (
	"bufio"
	"fmt"
	"io"
	"runtime"
	"time"

	"github.com/ntsd/cross-clipboard/pkg/clipboard"
	"github.com/ntsd/cross-clipboard/pkg/crypto"
	"github.com/ntsd/cross-clipboard/pkg/device"
)

// CreateReadData craete a new read streaming for host or peer
func (s *StreamHandler) CreateReadData(reader *bufio.Reader, dv *device.Device) {
	// sending this device data and public key
	pub, err := s.Config.PGPPrivateKey.GetPublicKey()
	if err != nil {
		s.ErrorChan <- fmt.Errorf("error to generate pubic key: %w", err)
		return
	}
	s.LogChan <- fmt.Sprintf("sending device info and public key to %s", dv.AddressInfo.ID.Pretty())
	deviceData, err := s.EncodeDeviceData(&DeviceData{
		Name:      s.Config.Username,
		Os:        runtime.GOOS,
		PublicKey: pub,
	})
	err = s.WriteData(dv.Writer, deviceData)
	if err != nil {
		s.ErrorChan <- fmt.Errorf("cannot send device data to %s: %w", dv.AddressInfo.ID.Pretty(), err)
		return
	}

	// loop for incoming message
	for {
		dataSize, err := readDataSize(reader)
		if err != nil {
			s.ErrorChan <- fmt.Errorf("error reading data size: %w", err)
			break
		}

		if dataSize <= 0 {
			s.ErrorChan <- fmt.Errorf("data size is less than 0: %d", dataSize)
			break
		}

		s.LogChan <- fmt.Sprintf("received data size %d", dataSize)

		buffer := make([]byte, dataSize)
		readBytes, err := io.ReadFull(reader, buffer)
		if err != nil {
			s.ErrorChan <- fmt.Errorf("error reading from buffer: %w", err)
			break
		}
		if readBytes != dataSize {
			s.ErrorChan <- fmt.Errorf("not reading full bytes read: %d size: %d", readBytes, dataSize)
			break
		}
		s.LogChan <- fmt.Sprintf("read data size %d", readBytes)

		clipboardData, deviceData, err := s.DecodeData(buffer)
		if err != nil {
			s.ErrorChan <- fmt.Errorf("error decoding data: %w", err)
			continue
		}

		if clipboardData != nil {
			s.LogChan <- fmt.Sprintf("received clipboard data, peer: %s size: %d", dv.AddressInfo.ID.Pretty(), clipboardData.DataSize)
			s.ClipboardManager.WriteClipboard(clipboard.Clipboard{
				IsImage: clipboardData.IsImage,
				Data:    clipboardData.Data,
				Size:    clipboardData.DataSize,
				Time:    time.UnixMicro(clipboardData.Time),
			})
		}

		if deviceData != nil {
			s.LogChan <- fmt.Sprintf("received device data, peer: %s", dv.AddressInfo.ID.Pretty())
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
	s.LogChan <- fmt.Sprintf("ending read stream for peer: %s", dv.AddressInfo.ID.Pretty())
}

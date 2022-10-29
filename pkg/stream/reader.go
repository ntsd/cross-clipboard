package stream

import (
	"bufio"
	"fmt"
	"io"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/ntsd/cross-clipboard/pkg/clipboard"
	"github.com/ntsd/cross-clipboard/pkg/device"
	"github.com/ntsd/cross-clipboard/pkg/xerror"
)

const limitDataSize = 1 << 20 // data size to avoid to read (100 MB)

// CreateReadData craete a new read streaming for host or peer
func (s *StreamHandler) CreateReadData(reader *bufio.Reader, dv *device.Device) {
	s.logChan <- fmt.Sprintf("sending device info and public key to %s", dv.AddressInfo.ID.Pretty())

	s.sendDeviceData(dv)

	// loop for incoming message
disconnect:
	for {
		dataSize, err := readDataSize(reader)
		if err != nil {
			if err == network.ErrReset { // error stream reset because it unusual stream end
				s.errorChan <- xerror.NewRuntimeErrorf("peer %s stream reset", dv.AddressInfo.ID.Pretty()).Wrap(err)
				dv.Status = device.StatusDisconnected
				s.deviceManager.UpdateDevice(dv)
				break disconnect
			}

			s.errorChan <- xerror.NewRuntimeError("error reading data size").Wrap(err)
			dv.Status = device.StatusError
			s.deviceManager.UpdateDevice(dv)
			break disconnect
		}

		if dataSize <= 0 {
			s.errorChan <- xerror.NewRuntimeErrorf("data size < 0, size %d", dataSize)
			dv.Status = device.StatusError
			s.deviceManager.UpdateDevice(dv)
			break disconnect
		}

		// avoid to read big data from stream
		if dataSize > limitDataSize {
			s.errorChan <- xerror.NewRuntimeErrorf("data size %d > limit data size %d", dataSize, limitDataSize)
			dv.Status = device.StatusBlocked
			s.deviceManager.UpdateDevice(dv)
			break disconnect
		}

		// skip clipboard size when data more than config max size
		if dataSize > s.config.MaxSize {
			s.errorChan <- xerror.NewRuntimeErrorf("data size %d > config max size %d", dataSize, s.config.MaxSize)
			reader.Discard(dataSize)
			continue
		}

		buffer := make([]byte, dataSize)
		readBytes, err := io.ReadFull(reader, buffer)
		if err != nil {
			s.errorChan <- xerror.NewRuntimeError("error reading from buffer").Wrap(err)
			dv.Status = device.StatusError
			s.deviceManager.UpdateDevice(dv)
			break disconnect
		}
		if readBytes != dataSize {
			s.errorChan <- xerror.NewRuntimeErrorf("not reading full bytes read: %d size: %d", readBytes, dataSize)
			dv.Status = device.StatusError
			s.deviceManager.UpdateDevice(dv)
			break disconnect
		}

		clipboardData, deviceData, signal, err := s.decodeData(buffer)
		if err != nil {
			s.errorChan <- xerror.NewRuntimeError("error decoding data").Wrap(err)
			dv.Status = device.StatusError
			s.deviceManager.UpdateDevice(dv)
			break disconnect
		}

		if signal != nil {
			s.logChan <- fmt.Sprintf("received signal %v, peer: %s", signal, dv.AddressInfo.ID.Pretty())
			switch *signal {
			case SignalDisconnect:
				dv.Status = device.StatusDisconnected
				s.deviceManager.UpdateDevice(dv)
				break disconnect
			case SignalRequestDeviceData:
				s.sendDeviceData(dv)
			}
		}

		if clipboardData != nil {
			s.clipboardManager.WriteClipboard(clipboard.FromProtobuf(clipboardData, dv))
			s.logChan <- fmt.Sprintf("received clipboard data, peer: %s size: %d", dv.AddressInfo.ID.Pretty(), clipboardData.DataSize)
		}

		if deviceData != nil {
			s.logChan <- fmt.Sprintf("received device data, peer: %s", dv.AddressInfo.ID.Pretty())

			s.logChan <- fmt.Sprintf("%s wanted to connect", deviceData.Name)
			dv.UpdateFromProtobuf(deviceData)

			if dv.PgpEncrypter == nil {
				dv.Status = device.StatusPending

				if s.config.AutoTrust {
					dv.Trust()
					s.logChan <- fmt.Sprintf("trusted %s by auto trust", deviceData.Name)
				}
			} else {
				dv.Status = device.StatusConnected
			}

			s.deviceManager.UpdateDevice(dv)
		}
	}

	s.logChan <- fmt.Sprintf("ending read stream for peer: %s", dv.AddressInfo.ID.Pretty())

	err := dv.Stream.Close()
	if err != nil {
		if err == network.ErrReset { // check stream already reset
			s.logChan <- fmt.Sprintf("peer %s stream already reset", dv.AddressInfo.ID.Pretty())
		}
		s.errorChan <- fmt.Errorf("can not close stream for peer %s: %w", dv.AddressInfo.ID, err)
	}
}

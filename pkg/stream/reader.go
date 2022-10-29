package stream

import (
	"bufio"
	"fmt"
	"io"

	"github.com/ntsd/cross-clipboard/pkg/clipboard"
	"github.com/ntsd/cross-clipboard/pkg/device"
	"github.com/ntsd/cross-clipboard/pkg/xerror"
)

// CreateReadData craete a new read streaming for host or peer
func (s *StreamHandler) CreateReadData(reader *bufio.Reader, dv *device.Device) {
	s.logChan <- fmt.Sprintf("sending device info and public key to %s", dv.AddressInfo.ID.Pretty())

	s.sendDeviceData(dv)

	// loop for incoming message
disconnect:
	for {
		dataSize, err := readDataSize(reader)
		if err != nil {
			s.errorChan <- xerror.NewRuntimeError("error reading data size").Wrap(err)
			dv.Status = device.StatusError
			s.deviceManager.UpdateDevice(dv)
			break
		}

		if dataSize <= 0 {
			s.errorChan <- xerror.NewRuntimeErrorf("data size < 0: %d", dataSize)
			dv.Status = device.StatusError
			s.deviceManager.UpdateDevice(dv)
			break
		}

		buffer := make([]byte, dataSize)
		readBytes, err := io.ReadFull(reader, buffer)
		if err != nil {
			s.errorChan <- xerror.NewRuntimeError("error reading from buffer").Wrap(err)
			dv.Status = device.StatusError
			s.deviceManager.UpdateDevice(dv)
			break
		}
		if readBytes != dataSize {
			s.errorChan <- xerror.NewRuntimeErrorf("not reading full bytes read: %d size: %d", readBytes, dataSize)
			dv.Status = device.StatusError
			s.deviceManager.UpdateDevice(dv)
			break
		}

		clipboardData, deviceData, signal, err := s.decodeData(buffer)
		if err != nil {
			s.errorChan <- xerror.NewRuntimeError("error decoding data").Wrap(err)
			dv.Status = device.StatusError
			s.deviceManager.UpdateDevice(dv)
			break
		}

		if signal != nil {
			s.logChan <- fmt.Sprintf("received signal %v, peer: %s size: %d", signal, dv.AddressInfo.ID.Pretty(), clipboardData.DataSize)
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
			s.logChan <- fmt.Sprintf("received clipboard data, peer: %s size: %d", dv.AddressInfo.ID.Pretty(), clipboardData.DataSize)
			s.clipboardManager.WriteClipboard(clipboard.FromProtobuf(clipboardData, dv))
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
		s.errorChan <- fmt.Errorf("can not close stream for peer %s: %w", dv.AddressInfo.ID, err)
	}
}

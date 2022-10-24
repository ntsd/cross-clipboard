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
exit:
	for {
		dataSize, err := readDataSize(reader)
		if err != nil {
			s.errorChan <- xerror.NewRuntimeError("error reading data size").Wrap(err)
			dv.Status = device.StatusError // TODO: handle peer exit to disconnect
			s.deviceManager.UpdateDevice(dv)
			break
		}

		if dataSize <= 0 {
			s.errorChan <- xerror.NewRuntimeErrorf("data size < 0: %d", dataSize)
			dv.Status = device.StatusError
			s.deviceManager.UpdateDevice(dv)
			break
		}

		s.logChan <- fmt.Sprintf("received data size %d", dataSize)

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
		s.logChan <- fmt.Sprintf("read data size %d", readBytes)

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
				break exit
			case SignalRequestDeviceData:
				s.sendDeviceData(dv)
			}
		}

		if clipboardData != nil {
			s.logChan <- fmt.Sprintf("received clipboard data, peer: %s size: %d", dv.AddressInfo.ID.Pretty(), clipboardData.DataSize)
			s.clipboardManager.WriteClipboard(clipboard.FromProtobuf(clipboardData))
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
}

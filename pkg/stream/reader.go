package stream

import (
	"bufio"
	"fmt"
	"io"
	"runtime"

	"github.com/ntsd/cross-clipboard/pkg/clipboard"
	"github.com/ntsd/cross-clipboard/pkg/device"
	"github.com/ntsd/cross-clipboard/pkg/protobuf"
	"github.com/ntsd/cross-clipboard/pkg/xerror"
)

// CreateReadData craete a new read streaming for host or peer
func (s *StreamHandler) CreateReadData(reader *bufio.Reader, dv *device.Device) {
	// sending this device data and public key
	pub, err := s.Config.PGPPrivateKey.GetPublicKey()
	if err != nil {
		s.ErrorChan <- xerror.NewFatalError("error to generate pubic key").Wrap(err)
		return
	}
	s.LogChan <- fmt.Sprintf("sending device info and public key to %s", dv.AddressInfo.ID.Pretty())
	deviceData, err := s.EncodeDeviceData(&protobuf.DeviceData{
		Name:      s.Config.Username,
		Os:        runtime.GOOS,
		PublicKey: pub,
	})
	err = s.WriteData(dv.Writer, deviceData)
	if err != nil {
		s.ErrorChan <- xerror.NewRuntimeErrorf("cannot send device data to %s", dv.AddressInfo.ID.Pretty()).Wrap(err)
		dv.Status = device.StatusError
		return
	}

	// loop for incoming message
	for {
		dataSize, err := readDataSize(reader)
		if err != nil {
			s.ErrorChan <- xerror.NewRuntimeError("error reading data size").Wrap(err)
			dv.Status = device.StatusError
			break
		}

		if dataSize <= 0 {
			s.ErrorChan <- xerror.NewRuntimeErrorf("data size < 0: %d", dataSize)
			dv.Status = device.StatusError
			break
		}

		s.LogChan <- fmt.Sprintf("received data size %d", dataSize)

		buffer := make([]byte, dataSize)
		readBytes, err := io.ReadFull(reader, buffer)
		if err != nil {
			s.ErrorChan <- xerror.NewRuntimeError("error reading from buffer").Wrap(err)
			dv.Status = device.StatusError
			break
		}
		if readBytes != dataSize {
			s.ErrorChan <- xerror.NewRuntimeErrorf("not reading full bytes read: %d size: %d", readBytes, dataSize)
			dv.Status = device.StatusError
			break
		}
		s.LogChan <- fmt.Sprintf("read data size %d", readBytes)

		clipboardData, deviceData, err := s.DecodeData(buffer)
		if err != nil {
			s.ErrorChan <- xerror.NewRuntimeError("error decoding data").Wrap(err)
			dv.Status = device.StatusError
			break
		}

		if clipboardData != nil {
			s.LogChan <- fmt.Sprintf("received clipboard data, peer: %s size: %d", dv.AddressInfo.ID.Pretty(), clipboardData.DataSize)
			s.ClipboardManager.WriteClipboard(clipboard.FromProtobuf(clipboardData))
		}

		if deviceData != nil {
			s.LogChan <- fmt.Sprintf("received device data, peer: %s", dv.AddressInfo.ID.Pretty())

			s.LogChan <- fmt.Sprintf("%s wanted to connect", deviceData.Name)
			dv.UpdateFromProtobuf(deviceData)
			dv.Status = device.StatusPending

			if s.Config.AutoTrust {
				dv.Status = device.StatusConnected
				s.LogChan <- fmt.Sprintf("trusted %s by auto trust", deviceData.Name)
			}

			s.DeviceManager.UpdateDevice(dv)
		}
	}
	s.LogChan <- fmt.Sprintf("ending read stream for peer: %s", dv.AddressInfo.ID.Pretty())
}

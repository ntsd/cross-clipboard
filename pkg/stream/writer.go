package stream

import (
	"bufio"
	"fmt"
	"runtime"
	"time"

	"github.com/ntsd/cross-clipboard/pkg/clipboard"
	"github.com/ntsd/cross-clipboard/pkg/device"
	"github.com/ntsd/cross-clipboard/pkg/protobuf"
	"github.com/ntsd/cross-clipboard/pkg/xerror"
)

// CreateWriteData handle clipboad channel and write to all peers and host
func (s *StreamHandler) CreateWriteData() {
	// waiting for clipboard data
loop:
	for {
		select {
		case textBytes, ok := <-s.clipboardManager.ReadTextChannel:
			if !ok {
				break loop
			}
			s.sendClipboard(textBytes, false)
		case imageBytes, ok := <-s.clipboardManager.ReadImageChannel:
			if !ok {
				break loop
			}
			s.sendClipboard(imageBytes, true)
		}
	}
	s.logChan <- "ended write streams"
}

func (s *StreamHandler) sendClipboard(clipboardBytes []byte, isImage bool) {
	clipboardLength := len(clipboardBytes)
	if clipboardLength == 0 {
		// ignore empty clipboard data
		s.logChan <- "the clipboard is empty, ignoring"
		return
	}

	cb := clipboard.Clipboard{
		IsImage: isImage,
		Data:    clipboardBytes,
		Size:    uint32(clipboardLength),
		Time:    time.Now(),
	}

	clipboardData := cb.ToProtobuf()

	// send data to each devices
	for name, dv := range s.deviceManager.Devices {
		if dv.Status == device.StatusPending {
			s.sendSignal(dv, SignalRequestDeviceData)
		}

		if dv.Status != device.StatusConnected {
			// skip disconnected devices
			continue
		}

		// avoid sending back to where it received
		if s.clipboardManager.IsReceivedClipboardFromDevice(dv) {
			continue
		}

		if dv.PgpEncrypter == nil {
			s.errorChan <- xerror.NewRuntimeErrorf("not found pgp encrypter for device %s", name)
			dv.Status = device.StatusError
			s.deviceManager.UpdateDevice(dv)
			// todo request for public key
			continue
		}

		s.logChan <- fmt.Sprintf("sending data to peer: %s len: %d", name, clipboardLength)

		clipboardDataBytes, err := s.encodeClipboardData(dv, clipboardData)
		if err != nil {
			s.errorChan <- xerror.NewRuntimeError("error encoding data").Wrap(err)
			dv.Status = device.StatusError
			s.deviceManager.UpdateDevice(dv)
			continue
		}

		err = s.writeData(dv.Writer, clipboardDataBytes)
		if err != nil {
			s.logChan <- fmt.Sprintf("error to send data for peer: %s", name)
			dv.Status = device.StatusError
			s.deviceManager.UpdateDevice(dv)
		}
	}

	s.clipboardManager.UpdateClipboard(cb)
}

// sendDeviceData send device data to the giving device
func (s *StreamHandler) sendDeviceData(dv *device.Device) {
	pub, err := s.config.PGPPrivateKey.GetPublicKey()
	if err != nil {
		s.errorChan <- xerror.NewFatalError("error to generate pubic key").Wrap(err)
		return
	}
	deviceData, err := s.encodeDeviceData(&protobuf.DeviceData{
		Name:      s.config.Username,
		Os:        runtime.GOOS,
		PublicKey: pub,
	})
	if err != nil {
		s.errorChan <- xerror.NewRuntimeError("cannot encode device data").Wrap(err)
		return
	}
	err = s.writeData(dv.Writer, deviceData)
	if err != nil {
		dv.Status = device.StatusError
		s.deviceManager.UpdateDevice(dv)
		s.errorChan <- xerror.NewRuntimeErrorf("cannot send device data to %s", dv.AddressInfo.ID.Pretty()).Wrap(err)
	}
}

// sendSignal send signal to device
func (s *StreamHandler) sendSignal(dv *device.Device, signal Signal) {
	signalData, err := s.encodeSignal(signal)
	if err != nil {
		s.errorChan <- xerror.NewRuntimeError("cannot encode signal").Wrap(err)
		return
	}
	err = s.writeData(dv.Writer, signalData)
	if err != nil {
		dv.Status = device.StatusError
		s.deviceManager.UpdateDevice(dv)
		s.errorChan <- xerror.NewRuntimeErrorf("cannot send signal to %s", dv.AddressInfo.ID.Pretty()).Wrap(err)
	}
}

// writeData write data to the writer
func (s *StreamHandler) writeData(w *bufio.Writer, data []byte) error {
	_, err := w.Write(data)
	if err != nil {
		s.errorChan <- xerror.NewRuntimeError("error writing to buffer").Wrap(err)
		return err
	}

	err = w.Flush()
	if err != nil {
		s.errorChan <- xerror.NewRuntimeError("error flushing buffer").Wrap(err)
		return err
	}
	return nil
}

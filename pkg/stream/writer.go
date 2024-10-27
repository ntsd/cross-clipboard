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
readClipboardLoop:
	for {
		select {
		case textBytes, ok := <-s.clipboardManager.ReadTextChannel:
			if !ok {
				break readClipboardLoop
			}
			s.sendClipboard(textBytes, false)
		case imageBytes, ok := <-s.clipboardManager.ReadImageChannel:
			if !ok {
				break readClipboardLoop
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

	if clipboardLength > s.config.MaxSize {
		s.errorChan <- xerror.NewRuntimeErrorf("clipboard size %d > config max size %d", clipboardLength, s.config.MaxSize)
		return
	}

	cb := &clipboard.Clipboard{
		IsImage: isImage,
		Data:    clipboardBytes,
		Size:    uint32(clipboardLength),
		Time:    time.Now(),
	}

	isReceivedClipboard := s.clipboardManager.IsReceivedClipboard(clipboardBytes)
	if !isReceivedClipboard { // the clipboard come from this device
		s.clipboardManager.AddClipboardToHistory(cb)
	}

	clipboardData := cb.ToProtobuf()

	// send data to each devices
	for name, dv := range s.deviceManager.Devices {
		if dv.Status == device.StatusPending {
			// request for public key
			s.SendSignal(dv, SignalRequestDeviceData)
			continue
		}

		if dv.Status != device.StatusConnected {
			// skip disconnected devices
			continue
		}

		// avoid sending back to where it received
		if isReceivedClipboard && s.clipboardManager.IsReceivedDevice(dv) {
			continue
		}

		if dv.PgpEncrypter == nil {
			s.errorChan <- xerror.NewRuntimeErrorf("not found pgp encrypter for device %s", name)
			dv.Status = device.StatusError
			s.deviceManager.UpdateDevice(dv)
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
		s.errorChan <- xerror.NewRuntimeErrorf("cannot send device data to %s", dv.AddressInfo.ID).Wrap(err)
	}
}

// SendSignal send signal to device
func (s *StreamHandler) SendSignal(dv *device.Device, signal Signal) {
	signalData, err := s.encodeSignal(signal)
	if err != nil {
		s.errorChan <- xerror.NewRuntimeError("cannot encode signal").Wrap(err)
		return
	}
	err = s.writeData(dv.Writer, signalData)
	if err != nil {
		dv.Status = device.StatusError
		s.deviceManager.UpdateDevice(dv)
		s.errorChan <- xerror.NewRuntimeErrorf("cannot send signal to %s", dv.AddressInfo.ID).Wrap(err)
	}
}

// writeData write data to the writer
func (s *StreamHandler) writeData(w *bufio.Writer, data []byte) error {
	_, err := w.Write(data)
	if err != nil {
		return xerror.NewRuntimeError("error writing to buffer").Wrap(err)
	}

	err = w.Flush()
	if err != nil {
		return xerror.NewRuntimeError("error flushing buffer").Wrap(err)
	}
	return nil
}

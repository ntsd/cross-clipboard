package stream

import (
	"bufio"
	"fmt"
	"time"

	"github.com/ntsd/cross-clipboard/pkg/clipboard"
	"github.com/ntsd/cross-clipboard/pkg/device"
	"github.com/ntsd/cross-clipboard/pkg/xerror"
)

// CreateWriteData handle clipboad channel and write to all peers and host
func (s *StreamHandler) CreateWriteData() {
	// waiting for clipboard data
loop:
	for {
		select {
		case textBytes, ok := <-s.ClipboardManager.ReadTextChannel:
			if !ok {
				break loop
			}
			s.sendClipboard(textBytes, false)
		case imageBytes, ok := <-s.ClipboardManager.ReadImageChannel:
			if !ok {
				break loop
			}
			s.sendClipboard(imageBytes, true)
		}
	}
	s.LogChan <- "ended write streams"
}

func (s *StreamHandler) sendClipboard(clipboardBytes []byte, isImage bool) {
	clipboardLength := len(clipboardBytes)
	if clipboardLength == 0 {
		// ignore empty clipboard data
		s.LogChan <- "the clipboard is empty, ignoring"
		return
	}

	now := time.Now()

	cb := clipboard.Clipboard{
		IsImage: isImage,
		Data:    clipboardBytes,
		Size:    uint32(clipboardLength),
		Time:    now,
	}

	// set current clipbaord to avoid recursive
	s.ClipboardManager.AddClipboard(cb)

	clipboardData := cb.ToProtobuf()

	// send data to each devices
	for name, dv := range s.DeviceManager.Devices {
		if dv.Status != device.StatusConnected {
			// skip disconnected devices
			continue
		}

		if dv.PgpEncrypter == nil {
			s.ErrorChan <- xerror.NewRuntimeErrorf("not found pgp encrypter for device %s", name)
			dv.Status = device.StatusError
			s.DeviceManager.UpdateDevice(dv)
			// todo request for public key
			continue
		}

		s.LogChan <- fmt.Sprintf("sending data to peer: %s len: %d", name, clipboardLength)

		clipboardDataBytes, err := s.EncodeClipboardData(dv, clipboardData)
		if err != nil {
			s.ErrorChan <- xerror.NewRuntimeError("error encoding data").Wrap(err)
			dv.Status = device.StatusError
			s.DeviceManager.UpdateDevice(dv)
			continue
		}

		err = s.WriteData(dv.Writer, clipboardDataBytes)
		if err != nil {
			s.LogChan <- fmt.Sprintf("error to send data for peer: %s", name)
			dv.Status = device.StatusError
			s.DeviceManager.UpdateDevice(dv)
		}
	}
}

// WriteData write data to the writer
func (s *StreamHandler) WriteData(w *bufio.Writer, data []byte) error {
	_, err := w.Write(data)
	if err != nil {
		s.ErrorChan <- xerror.NewRuntimeError("error writing to buffer").Wrap(err)
		return err
	}

	err = w.Flush()
	if err != nil {
		s.ErrorChan <- xerror.NewRuntimeError("error flushing buffer").Wrap(err)
		return err
	}
	return nil
}

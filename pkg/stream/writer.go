package stream

import (
	"bufio"
	"fmt"
	"time"

	"github.com/ntsd/cross-clipboard/pkg/clipboard"
	"github.com/ntsd/cross-clipboard/pkg/device"
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
			err := s.sendClipboard(textBytes, false)
			if err != nil {
				s.ErrorChan <- fmt.Errorf("error sending text clipboard data: %w", err)
				break loop
			}
		case imageBytes, ok := <-s.ClipboardManager.ReadImageChannel:
			if !ok {
				break loop
			}
			err := s.sendClipboard(imageBytes, true)
			if err != nil {
				s.ErrorChan <- fmt.Errorf("error sending image clipboard data: %w", err)
				break loop
			}
		}
	}
	s.LogChan <- fmt.Sprintf("ending write streams")
}

func (s *StreamHandler) sendClipboard(clipboardBytes []byte, isImage bool) error {
	clipboardLength := len(clipboardBytes)
	if clipboardLength == 0 {
		// ignore empty clipboard data
		s.LogChan <- "the clipboard is empty, ignoring"
		return nil
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
			s.ErrorChan <- fmt.Errorf("device %s status is not connected", name)
			continue
		}

		if dv.PgpEncrypter == nil {
			s.ErrorChan <- fmt.Errorf("not found pgp encrypter for device %s", name)
			dv.Status = device.StatusError
			// todo request for public key
			continue
		}

		s.LogChan <- fmt.Sprintf("sending data to peer: %s len: %d", name, clipboardLength)

		clipboardDataBytes, err := s.EncodeClipboardData(dv, clipboardData)
		if err != nil {
			s.ErrorChan <- fmt.Errorf("error encoding data: %w", err)
			dv.Status = device.StatusError
			continue
		}

		err = s.WriteData(dv.Writer, clipboardDataBytes)
		if err != nil {
			s.LogChan <- fmt.Sprintf("error to send data for peer: %s", name)
			dv.Status = device.StatusError
		}
	}

	return nil
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

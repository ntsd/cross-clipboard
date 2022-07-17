package stream

import (
	"bufio"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/ntsd/cross-clipboard/pkg/clipboard"
	"github.com/ntsd/cross-clipboard/pkg/devicemanager"
	"google.golang.org/protobuf/proto"
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
	ClipboardManager *clipboard.ClipboardManager
	DeviceManager    *devicemanager.DeviceManager
	HostReader       *bufio.Reader
	HostWriter       *bufio.Writer
	LogChan          chan string
	ErrorChan        chan error
}

// NewStreamHandler initial new stream handler
func NewStreamHandler(
	cp *clipboard.ClipboardManager,
	deviceManager *devicemanager.DeviceManager,
	logChan chan string,
	errorChan chan error,
) *StreamHandler {
	s := &StreamHandler{
		ClipboardManager: cp,
		DeviceManager:    deviceManager,
		LogChan:          logChan,
		ErrorChan:        errorChan,
	}
	go s.CreateWriteData()
	return s
}

// HandleStream handler when a peer connect this host
func (s *StreamHandler) HandleStream(stream network.Stream) {
	s.LogChan <- fmt.Sprintf("got a new stream from %s", stream.Conn().RemotePeer())

	// Create a new peer
	s.HostReader = bufio.NewReader(stream)
	s.HostWriter = bufio.NewWriter(stream)
	go s.CreateReadData(s.HostReader, "host")

	// 'stream' will stay open until you close it (or the other side closes it).
}

// CreateReadData craete a new read streaming for host or peer
func (s *StreamHandler) CreateReadData(reader *bufio.Reader, name string) {
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
			s.LogChan <- fmt.Sprintf("received clipboard data, peer: %s size: %d", name, clipboardData.DataSize)
			s.ClipboardManager.WriteClipboard(clipboard.Clipboard{
				IsImage: clipboardData.IsImage,
				Data:    clipboardData.Data,
				Size:    clipboardData.DataSize,
				Time:    time.UnixMicro(clipboardData.Time),
			})
		}

		if deviceData != nil {
			s.LogChan <- fmt.Sprintf("received device data, peer: %s", name)
		}
	}
	s.LogChan <- fmt.Sprintf("ending read stream for peer: %s", name)
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

		clipboardDataBytes, err := s.EncodeClipboardData(&ClipboardData{
			IsImage:  false,
			Data:     clipboardBytes,
			DataSize: uint32(length),
			Time:     now.Unix(),
		})
		if err != nil {
			s.ErrorChan <- fmt.Errorf("error encoding data: %w", err)
			continue
		}

		// send data to each devices
		for name, d := range s.DeviceManager.Devices {
			s.LogChan <- fmt.Sprintf("sending data to peer: %s size: %d data: %s", name, length, string(clipboardBytes))
			err := s.WriteData(d.Writer, clipboardDataBytes)
			if err != nil {
				s.LogChan <- fmt.Sprintf("ending write stream %s", name)
				s.DeviceManager.RemoveDevice(d)
			}
		}

		// send data to host
		if s.HostWriter != nil {
			s.LogChan <- fmt.Sprintf("sending data to host size: %d data: %s", length, string(clipboardBytes))
			s.WriteData(s.HostWriter, clipboardDataBytes)
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

func (s *StreamHandler) EncodeClipboardData(data *ClipboardData) ([]byte, error) {
	// create proto clipboard data
	clipboardDataBytes, err := proto.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshaling clipboard data: %w", err)
	}

	// append DATA TYPE
	clipboardDataBytes = append(clipboardDataBytes, DATA_TYPE_CLIPBOARD)

	// append EOF byte
	clipboardDataBytes = append(clipboardDataBytes, EOF)

	return clipboardDataBytes, nil
}

func (s *StreamHandler) EncodeDeviceData(data *DeviceData) ([]byte, error) {
	// create proto clipboard data
	clipboardDataBytes, err := proto.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshaling device data: %w", err)
	}

	// append DATA TYPE
	clipboardDataBytes = append(clipboardDataBytes, DATA_TYPE_DEVICE)

	// append EOF byte
	clipboardDataBytes = append(clipboardDataBytes, EOF)

	return clipboardDataBytes, nil
}

func (s *StreamHandler) DecodeData(bytes []byte) (*ClipboardData, *DeviceData, error) {
	length := len(bytes) - 2
	if length <= 0 {
		return nil, nil, fmt.Errorf("error decoding data: data length is 0")
	}

	// get data type from the last byte before EOF
	dataType := bytes[length-1]

	// remove EOF and data type bytes
	bytes = bytes[:length-1]

	switch dataType {
	case DATA_TYPE_CLIPBOARD:
		clipboardData := &ClipboardData{}
		err := proto.Unmarshal(bytes, clipboardData)
		if err != nil {
			return nil, nil, fmt.Errorf("error unmarshaling clipboard data: %w", err)
		}
		return clipboardData, nil, nil
	case DATA_TYPE_DEVICE:
		deviceData := &DeviceData{}
		err := proto.Unmarshal(bytes, deviceData)
		if err != nil {
			return nil, nil, fmt.Errorf("error unmarshaling device data: %w", err)
		}
		return nil, deviceData, nil
	default:
		return nil, nil, fmt.Errorf("error decoding data: unknown data type")
	}
}

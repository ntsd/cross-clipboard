package stream

import (
	"bufio"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/ntsd/cross-clipboard/pkg/clipboard"
	"github.com/ntsd/cross-clipboard/pkg/p2p"
	"google.golang.org/protobuf/proto"
)

const EOF byte = 0x00

// StreamHandler struct for stream handler
type StreamHandler struct {
	ClipboardManager *clipboard.ClipboardManager
	HostReader       *bufio.Reader
	HostWriter       *bufio.Writer
	Peers            map[string]*p2p.Peer
	LogChan          chan string
	ErrorChan        chan error
}

// NewStreamHandler initial new stream handler
func NewStreamHandler(cp *clipboard.ClipboardManager, logChan chan string, errorChan chan error, peers map[string]*p2p.Peer) *StreamHandler {
	s := &StreamHandler{
		ClipboardManager: cp,
		Peers:            peers,
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

		length := len(bytes) - 1
		if length > 0 {
			// remove EOF bytes
			bytes = bytes[:length]

			clipbaodData := &ClipboardData{}
			err = proto.Unmarshal(bytes, clipbaodData)
			if err != nil {
				s.ErrorChan <- fmt.Errorf("error unmarshaling data: %w", err)
				continue
			}

			s.LogChan <- fmt.Sprintf("received data from peer: %s size: %d data: %s", name, length, string(bytes))
			s.ClipboardManager.WriteClipboard(clipboard.Clipboard{
				IsImage: clipbaodData.IsImage,
				Data:    clipbaodData.Data,
				Size:    clipbaodData.DataSize,
				Time:    time.Unix(clipbaodData.Time, 0),
			})
		}
	}
	s.LogChan <- fmt.Sprintf("ending read stream for peer: %s", name)
}

// CreateWriteData handle clipboad channel and write to all peers and host
func (s *StreamHandler) CreateWriteData() {
	for clipboardBytes := range s.ClipboardManager.ReadChannel {
		length := len(clipboardBytes)
		if length > 0 {
			now := time.Now()

			// set current clipbaord to avoid recursive
			s.ClipboardManager.AddClipboard(clipboard.Clipboard{
				Data: clipboardBytes,
				Size: uint32(length),
				Time: now,
			})

			// create proto clipboard data
			clipboardDataBytes, err := proto.Marshal(&ClipboardData{
				IsImage:  false,
				Data:     clipboardBytes,
				DataSize: uint32(length),
				Time:     now.Unix(),
			})
			if err != nil {
				s.ErrorChan <- fmt.Errorf("ending write stream %w", err)
				continue
			}

			// append EOF byte
			clipboardDataBytes = append(clipboardDataBytes, EOF)

			for name, p := range s.Peers {
				s.LogChan <- fmt.Sprintf("sending data to peer: %s size: %d data: %s", name, length, string(clipboardBytes))
				err := s.WriteData(p.Writer, clipboardDataBytes)
				if err != nil {
					s.LogChan <- fmt.Sprintf("ending write stream %s", name)
					delete(s.Peers, name)
				}
			}

			if s.HostWriter != nil {
				s.LogChan <- fmt.Sprintf("sending data to host size: %d data: %s", length, string(clipboardBytes))
				s.WriteData(s.HostWriter, clipboardDataBytes)
			}
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

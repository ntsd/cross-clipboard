package stream

import (
	"bufio"
	"fmt"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/ntsd/cross-clipboard/pkg/clipboard"
	"github.com/ntsd/cross-clipboard/pkg/p2p"
)

const EOF byte = 0x00

type StreamHandler struct {
	Clipboard  *clipboard.Clipboard
	HostReader *bufio.Reader
	HostWriter *bufio.Writer
	Peers      map[string]*p2p.Peer
	LogChan    chan string
	ErrorChan  chan error
}

func NewStreamHandler(cp *clipboard.Clipboard, logChan chan string, errorChan chan error, peers map[string]*p2p.Peer) *StreamHandler {
	s := &StreamHandler{
		Clipboard: cp,
		Peers:     peers,
		LogChan:   logChan,
		ErrorChan: errorChan,
	}
	go s.CreateWriteData()
	return s
}

func (s *StreamHandler) HandleStream(stream network.Stream) {
	s.LogChan <- fmt.Sprintf("got a new stream from %s", stream.Conn().RemotePeer())

	// Create a buffer stream for non blocking read and write.
	// rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	// Create a new peer
	s.HostReader = bufio.NewReader(stream)
	s.HostWriter = bufio.NewWriter(stream)
	go s.CreateReadData(s.HostReader, "host")

	// 'stream' will stay open until you close it (or the other side closes it).
}

func (s *StreamHandler) CreateReadData(reader *bufio.Reader, name string) {
	for {
		bytes, err := reader.ReadBytes(EOF)
		if err != nil {
			s.ErrorChan <- fmt.Errorf("error reading from buffer: %w", err)
			break
		}
		// remove last bytes
		length := len(bytes) - 1
		if length > 0 {
			bytes = bytes[:length]
			s.LogChan <- fmt.Sprintf("received data from peer: %s \n size: %d data: %s\n", name, length, string(bytes))
			s.Clipboard.Write(bytes)
		}
	}
	s.LogChan <- fmt.Sprintf("ending read stream for peer: %s", name)
}

func (s *StreamHandler) CreateWriteData() {
	for clipboardBytes := range s.Clipboard.ReadChannel {
		length := len(clipboardBytes)
		if length > 0 {
			// set current clipbaord to avoid recursion
			s.Clipboard.CurrentClipboard = clipboardBytes

			// append EOF
			clipboardBytes = append(clipboardBytes, EOF)

			for name, p := range s.Peers {
				s.LogChan <- fmt.Sprintf("sending data to peer: %s \n size: %d data: %s\n", name, length, string(clipboardBytes))
				err := s.WriteData(p.Writer, clipboardBytes)
				if err != nil {
					delete(s.Peers, name)
				}
			}

			s.WriteData(s.HostWriter, clipboardBytes)
		}
	}
	s.LogChan <- fmt.Sprintf("ending write streams")
}

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

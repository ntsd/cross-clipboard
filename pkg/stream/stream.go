package stream

import (
	"bufio"
	"fmt"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/ntsd/cross-clipboard/pkg/clipboard"
)

const EOF byte = 0x00

type StreamHandler struct {
	Clipboard *clipboard.Clipboard
	Writers   map[string]*bufio.Writer
}

func NewStreamHandler(cp *clipboard.Clipboard) *StreamHandler {
	s := &StreamHandler{
		Clipboard: cp,
		Writers:   make(map[string]*bufio.Writer),
	}
	go s.CreateWriteData()
	return s
}

func (s *StreamHandler) HandleStream(stream network.Stream) {
	fmt.Println("Got a new stream!")

	// Create a buffer stream for non blocking read and write.
	// rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	go s.CreateReadData(bufio.NewReader(stream), "host")
	s.AddWriter(bufio.NewWriter(stream), "host")

	// 'stream' will stay open until you close it (or the other side closes it).
}

func (s *StreamHandler) CreateReadData(reader *bufio.Reader, name string) {
	for {
		bytes, err := reader.ReadBytes(EOF)
		if err != nil {
			fmt.Println("Error reading from buffer:", err)
			break
		}
		// remove last bytes
		length := len(bytes) - 1
		if length > 0 {
			bytes = bytes[:length]
			fmt.Printf("Received data from peer: %s \n size: %d data: %s\n", name, length, string(bytes))
			s.Clipboard.Write(bytes)
		}
	}
	fmt.Println("Ending read stream for peer:", name)
}

func (s *StreamHandler) CreateWriteData() {
	for clipboardBytes := range s.Clipboard.ReadChannel {
		length := len(clipboardBytes)
		if length > 0 {
			// set current clipbaord to avoid recursion
			s.Clipboard.CurrentClipboard = clipboardBytes

			// append EOF
			clipboardBytes = append(clipboardBytes, EOF)

			for name, writer := range s.Writers {
				fmt.Printf("Sending data to peer %s \n size: %d data: %s\n", name, length, string(clipboardBytes))

				_, err := writer.Write(clipboardBytes)
				if err != nil {
					fmt.Println("Error writing buffer:", err)
					delete(s.Writers, name)
				}

				err = writer.Flush()
				if err != nil {
					fmt.Println("Error flush writing:", err)
					delete(s.Writers, name)
				}
			}
		}
	}
	fmt.Println("Ending write streams")
}

func (s *StreamHandler) AddWriter(writer *bufio.Writer, name string) {
	s.Writers[name] = writer
}

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
}

func NewStreamHandler(cp *clipboard.Clipboard) *StreamHandler {
	return &StreamHandler{
		Clipboard: cp,
	}
}

func (s *StreamHandler) HandleStream(stream network.Stream) {
	fmt.Println("Got a new stream!")

	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	go s.ReadData(rw)
	go s.WriteData(rw)

	// 'stream' will stay open until you close it (or the other side closes it).
}

func (s *StreamHandler) ReadData(rw *bufio.ReadWriter) {
	for {
		bytes, err := rw.ReadBytes(EOF)
		if err != nil {
			fmt.Println("Error reading from buffer:", err)
			break
		}
		// remove last bytes
		length := len(bytes) - 1
		if length > 0 {
			bytes = bytes[:length]
			fmt.Printf("Received data from peer\n size: %d data: %s\n", length, string(bytes))
			s.Clipboard.Write(bytes)
		}
	}
	fmt.Println("Ending read stream")
}

func (s *StreamHandler) WriteData(rw *bufio.ReadWriter) {
	for clipboardBytes := range s.Clipboard.ReadChannel {
		length := len(clipboardBytes)
		if length > 0 {
			fmt.Printf("Sending data to peer\n size: %d data: %s\n", length, string(clipboardBytes))

			// append EOF
			clipboardBytes = append(clipboardBytes, EOF)

			_, err := rw.Write(clipboardBytes)
			if err != nil {
				fmt.Println("Error writing buffer:", err)
				break
			}

			err = rw.Flush()
			if err != nil {
				fmt.Println("Error flush writing:", err)
				break
			}
		}
	}
	fmt.Println("Ending write stream")
}

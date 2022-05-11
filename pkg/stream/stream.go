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
	// for {
	// 	str, err := rw.ReadString('\n')
	// 	if err != nil {
	// 		fmt.Println("Error reading from buffer:", err)
	// 	}

	// 	if str == "" {
	// 		return
	// 	}
	// 	if str != "\n" {
	// 		// Green console colour: 	\x1b[32m
	// 		// Reset console colour: 	\x1b[0m
	// 		fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
	// 	}
	// }
	for {
		bytes, err := rw.ReadBytes(EOF)
		if err != nil {
			fmt.Println("Error reading from buffer:", err)
		}
		fmt.Println("Received data from peer:", string(bytes))
		s.Clipboard.Write(bytes)
	}
}

func (s *StreamHandler) WriteData(rw *bufio.ReadWriter) {
	for clipboardBytes := range s.Clipboard.ReadChannel {
		fmt.Println("Sending data to peer:", string(clipboardBytes))
		rw.Write(clipboardBytes)
		rw.WriteByte(EOF)
		rw.Flush()
	}

	// stdReader := bufio.NewReader(os.Stdin)

	// for {
	// 	fmt.Print("> ")
	// 	sendData, err := stdReader.ReadString('\n')
	// 	if err != nil {
	// 		fmt.Println("Error reading from stdin")
	// 		panic(err)
	// 	}

	// 	_, err = rw.WriteString(fmt.Sprintf("%s\n", sendData))
	// 	if err != nil {
	// 		fmt.Println("Error writing to buffer", err)
	// 	}
	// 	err = rw.Flush()
	// 	if err != nil {
	// 		fmt.Println("Error flushing buffer :", err)
	// 	}
	// }
}

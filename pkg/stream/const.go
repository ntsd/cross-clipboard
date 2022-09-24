package stream

import "github.com/libp2p/go-libp2p-core/protocol"

const (
	PROTOCAL_ID protocol.ID = protocol.ID("/cross-clipboard/0.0.1")

	// DATA TYPE first byte of the message
	DATA_TYPE_DEVICE           byte = 0xFF
	DATA_TYPE_SECURE_CLIPBOARD byte = 0xFE
	DATA_TYPE_CLIPBOARD        byte = 0xFD
)

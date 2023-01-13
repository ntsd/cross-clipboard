package stream

import "github.com/libp2p/go-libp2p/core/protocol"

type (
	DataType byte
	Signal   byte
)

const (
	PROTOCAL_ID protocol.ID = protocol.ID("/cross-clipboard/0.0.1")

	// data type is the first byte after data size to identify the message type
	DataTypeDevice    DataType = 0xFF // use for device data
	DataTypeClipboard DataType = 0xFE // use for clipboard data

	// signal is the first byte after data size to identify the signal type
	SignalDisconnect        Signal = 0xFD // ending exit signal
	SignalRequestDeviceData Signal = 0xFC // request device data signal
)

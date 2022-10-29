package stream

import (
	"github.com/ntsd/cross-clipboard/pkg/protobuf"
	"github.com/ntsd/cross-clipboard/pkg/xerror"
	"google.golang.org/protobuf/proto"
)

// decodeData decode message data to protobuf type `| data size (int 4 bytes) | data type (enum 1 byte) | protobuf message (struct n bytes) |`
func (s *StreamHandler) decodeData(bytes []byte) (*protobuf.ClipboardData, *protobuf.DeviceData, *Signal, error) {
	length := len(bytes)
	if length <= 0 {
		return nil, nil, nil, xerror.NewRuntimeErrorf("data length <= 0: %d", length)
	}

	// get data type from the last byte before EOF
	dataType := bytes[0]

	// remove data type bytes
	bytes = bytes[1:]

	switch dataType {
	case byte(DataTypeClipboard):
		// decrypt clipboard data
		decrypedData, err := s.pgpDecrypter.DecryptMessage(bytes)
		if err != nil {
			return nil, nil, nil, xerror.NewRuntimeError("error to decrypt clipboard data").Wrap(err)
		}

		clipboardData := &protobuf.ClipboardData{}
		err = proto.Unmarshal(decrypedData, clipboardData)
		if err != nil {
			return nil, nil, nil, xerror.NewRuntimeError("error unmarshaling clipboard data").Wrap(err)
		}
		return clipboardData, nil, nil, nil
	case byte(DataTypeDevice):
		deviceData := &protobuf.DeviceData{}
		err := proto.Unmarshal(bytes, deviceData)
		if err != nil {
			return nil, nil, nil, xerror.NewRuntimeError("error unmarshaling device data").Wrap(err)
		}
		return nil, deviceData, nil, nil
	case byte(SignalDisconnect):
		sn := SignalDisconnect
		return nil, nil, &sn, nil
	case byte(SignalRequestDeviceData):
		sn := SignalRequestDeviceData
		return nil, nil, &sn, nil
	default:
		return nil, nil, nil, xerror.NewRuntimeError("unknown data type")
	}
}

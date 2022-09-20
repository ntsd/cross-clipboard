package stream

import (
	"github.com/ntsd/cross-clipboard/pkg/protobuf"
	"github.com/ntsd/cross-clipboard/pkg/xerror"
	"google.golang.org/protobuf/proto"
)

// DecodeData decode message data to protobuf type `| data size (int 4 bytes) | data type (enum 1 byte) | protobuf message (struct n bytes) |`
func (s *StreamHandler) DecodeData(bytes []byte) (*protobuf.ClipboardData, *protobuf.DeviceData, error) {
	length := len(bytes)
	if length <= 1 {
		return nil, nil, xerror.NewRuntimeErrorf("data length <= 0: %d", length)
	}

	// get data type from the last byte before EOF
	dataType := bytes[0]

	// remove data type bytes
	bytes = bytes[1:]

	switch dataType {
	case DATA_TYPE_CLIPBOARD, DATA_TYPE_SECURE_CLIPBOARD:
		// decrypt clipboard data
		if dataType == DATA_TYPE_SECURE_CLIPBOARD {
			decryped, err := s.pgpDecrypter.DecryptMessage(bytes)
			if err != nil {
				return nil, nil, xerror.NewRuntimeError("error to decrypt clipboard data").Wrap(err)
			}
			bytes = decryped
		}

		clipboardData := &protobuf.ClipboardData{}
		err := proto.Unmarshal(bytes, clipboardData)
		if err != nil {
			return nil, nil, xerror.NewRuntimeError("error unmarshaling clipboard data").Wrap(err)
		}
		return clipboardData, nil, nil
	case DATA_TYPE_DEVICE:
		deviceData := &protobuf.DeviceData{}
		err := proto.Unmarshal(bytes, deviceData)
		if err != nil {
			return nil, nil, xerror.NewRuntimeError("error unmarshaling device data").Wrap(err)
		}
		return nil, deviceData, nil
	default:
		return nil, nil, xerror.NewRuntimeError("unknown data type")
	}
}

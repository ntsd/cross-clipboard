package stream

import (
	"fmt"

	"github.com/ntsd/cross-clipboard/pkg/protobuf"
	"google.golang.org/protobuf/proto"
)

func (s *StreamHandler) DecodeData(bytes []byte) (*protobuf.ClipboardData, *protobuf.DeviceData, error) {
	length := len(bytes)
	if length <= 1 {
		return nil, nil, fmt.Errorf("error decoding data: data length <= 0")
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
				return nil, nil, fmt.Errorf("error to decrypt clipboard data: %w", err)
			}
			bytes = decryped
		}

		clipboardData := &protobuf.ClipboardData{}
		err := proto.Unmarshal(bytes, clipboardData)
		if err != nil {
			return nil, nil, fmt.Errorf("error unmarshaling clipboard data: %w", err)
		}
		return clipboardData, nil, nil
	case DATA_TYPE_DEVICE:
		deviceData := &protobuf.DeviceData{}
		err := proto.Unmarshal(bytes, deviceData)
		if err != nil {
			return nil, nil, fmt.Errorf("error unmarshaling device data: %w", err)
		}
		return nil, deviceData, nil
	default:
		return nil, nil, fmt.Errorf("error decoding data: unknown data type")
	}
}

package stream

import (
	"fmt"

	"google.golang.org/protobuf/proto"
)

func (s *StreamHandler) EncodeClipboardData(data *ClipboardData) ([]byte, error) {
	// create proto clipboard data
	clipboardDataBytes, err := proto.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshaling clipboard data: %w", err)
	}

	// append DATA TYPE
	clipboardDataBytes = append(clipboardDataBytes, DATA_TYPE_CLIPBOARD)

	// append EOF byte
	clipboardDataBytes = append(clipboardDataBytes, EOF)

	return clipboardDataBytes, nil
}

func (s *StreamHandler) EncodeDeviceData(data *DeviceData) ([]byte, error) {
	// create proto clipboard data
	clipboardDataBytes, err := proto.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshaling device data: %w", err)
	}

	// append DATA TYPE
	clipboardDataBytes = append(clipboardDataBytes, DATA_TYPE_DEVICE)

	// append EOF byte
	clipboardDataBytes = append(clipboardDataBytes, EOF)

	return clipboardDataBytes, nil
}

func (s *StreamHandler) DecodeData(bytes []byte) (*ClipboardData, *DeviceData, error) {
	length := len(bytes) - 2
	if length <= 0 {
		return nil, nil, fmt.Errorf("error decoding data: data length <= 0")
	}

	// get data type from the last byte before EOF
	dataType := bytes[length]

	// remove EOF and data type bytes
	bytes = bytes[:length]

	switch dataType {
	case DATA_TYPE_CLIPBOARD:
		clipboardData := &ClipboardData{}
		err := proto.Unmarshal(bytes, clipboardData)
		if err != nil {
			return nil, nil, fmt.Errorf("error unmarshaling clipboard data: %w", err)
		}
		return clipboardData, nil, nil
	case DATA_TYPE_DEVICE:
		deviceData := &DeviceData{}
		err := proto.Unmarshal(bytes, deviceData)
		if err != nil {
			return nil, nil, fmt.Errorf("error unmarshaling device data: %w", err)
		}
		return nil, deviceData, nil
	default:
		return nil, nil, fmt.Errorf("error decoding data: unknown data type")
	}
}

package stream

import (
	"fmt"

	"github.com/ntsd/cross-clipboard/pkg/device"
	"google.golang.org/protobuf/proto"
)

const (
	// DATA TYPE first byte of the message
	DATA_TYPE_DEVICE           byte = 0xFF
	DATA_TYPE_SECURE_CLIPBOARD byte = 0xFE
	DATA_TYPE_CLIPBOARD        byte = 0xFD
)

// EncodeClipboardData encode data for stream package | size(bytes) int 4 bytes | data type 1 byte | message n bytes |
func (s *StreamHandler) EncodeClipboardData(dv *device.Device, clipboardData *ClipboardData) ([]byte, error) {
	packageData := []byte{}

	// create proto clipboard data
	clipboardDataBytes, err := proto.Marshal(clipboardData)
	if err != nil {
		return nil, fmt.Errorf("error marshaling clipboard data: %w", err)
	}
	dataSize := len(clipboardDataBytes)
	dataType := DATA_TYPE_CLIPBOARD

	// encrypt clipboard data
	if dataSize > 1024 {
		clipboardDataEncrypted, err := dv.PgpEncrypter.EncryptMessage(clipboardDataBytes)
		if err != nil {
			return nil, fmt.Errorf("error to encrypt clipboard data: %w", err)
		}
		dataSize = len(clipboardDataEncrypted)
		clipboardDataBytes = clipboardDataEncrypted
		dataType = DATA_TYPE_SECURE_CLIPBOARD
		s.LogChan <- fmt.Sprintf("dataSize: %d encrypted dataSize: %d", len(clipboardDataBytes), dataSize)
	}

	// append data size + 1 bytes for data type
	packageData = append(packageData, intToBytes(dataSize+1)...)
	// append DATA TYPE
	packageData = append(packageData, dataType)
	// append message
	packageData = append(packageData, clipboardDataBytes...)

	return packageData, nil
}

// EncodeDeviceData encode data for stream package | size(bytes) int 4 bytes | data type 1 byte | message n bytes |
func (s *StreamHandler) EncodeDeviceData(data *DeviceData) ([]byte, error) {
	packageData := []byte{}

	// create proto device data
	deviceDataBytes, err := proto.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshaling device data: %w", err)
	}
	dataSize := len(deviceDataBytes)

	// append data size + 1 bytes for data type
	packageData = append(packageData, intToBytes(dataSize+1)...)

	// append DATA TYPE
	packageData = append(packageData, DATA_TYPE_DEVICE)

	// append message
	packageData = append(packageData, deviceDataBytes...)

	return packageData, nil
}

func (s *StreamHandler) DecodeData(bytes []byte) (*ClipboardData, *DeviceData, error) {
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

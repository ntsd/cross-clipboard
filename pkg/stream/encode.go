package stream

import (
	"fmt"

	"github.com/ntsd/cross-clipboard/pkg/device"
	"github.com/ntsd/cross-clipboard/pkg/protobuf"
	"google.golang.org/protobuf/proto"
)

// EncodeClipboardData encode data for stream package | size(bytes) int 4 bytes | data type 1 byte | message n bytes |
func (s *StreamHandler) EncodeClipboardData(dv *device.Device, clipboardData *protobuf.ClipboardData) ([]byte, error) {
	packageData := []byte{}

	// create proto clipboard data
	clipboardDataBytes, err := proto.Marshal(clipboardData)
	if err != nil {
		return nil, fmt.Errorf("error marshaling clipboard data: %w", err)
	}
	dataSize := len(clipboardDataBytes)
	dataType := DATA_TYPE_CLIPBOARD

	// encrypt clipboard data
	if s.Config.IsEncryptEnabled {
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
func (s *StreamHandler) EncodeDeviceData(data *protobuf.DeviceData) ([]byte, error) {
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

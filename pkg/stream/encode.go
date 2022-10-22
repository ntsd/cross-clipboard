package stream

import (
	"fmt"

	"github.com/ntsd/cross-clipboard/pkg/device"
	"github.com/ntsd/cross-clipboard/pkg/protobuf"
	"github.com/ntsd/cross-clipboard/pkg/xerror"
	"google.golang.org/protobuf/proto"
)

// EncodeClipboardData encode data for stream package `| data size (int 4 bytes) | data type (enum 1 byte) | protobuf message (struct n bytes) |`
func (s *StreamHandler) EncodeClipboardData(dv *device.Device, clipboardData *protobuf.ClipboardData) ([]byte, error) {
	packageData := []byte{}

	// create proto clipboard data
	clipboardDataBytes, err := proto.Marshal(clipboardData)
	if err != nil {
		return nil, xerror.NewRuntimeError("error marshaling clipboard data").Wrap(err)
	}
	dataSize := len(clipboardDataBytes)
	dataType := DATA_TYPE_CLIPBOARD

	// encrypt clipboard data
	clipboardDataEncrypted, err := dv.PgpEncrypter.EncryptMessage(clipboardDataBytes)
	if err != nil {
		return nil, xerror.NewRuntimeError("error to encrypt clipboard data").Wrap(err)
	}
	encryptedDataSize := len(clipboardDataEncrypted)
	s.LogChan <- fmt.Sprintf("data size: %d encrypted data size: %d", dataSize, encryptedDataSize)

	// append data size + 1 bytes for data type
	packageData = append(packageData, intToBytes(encryptedDataSize+1)...)
	// append data type
	packageData = append(packageData, dataType)
	// append message
	packageData = append(packageData, clipboardDataEncrypted...)

	return packageData, nil
}

// EncodeDeviceData encode data for stream package `| data size (int 4 bytes) | data type (enum 1 byte) | protobuf message (struct n bytes) |`
func (s *StreamHandler) EncodeDeviceData(data *protobuf.DeviceData) ([]byte, error) {
	packageData := []byte{}

	// create proto device data
	deviceDataBytes, err := proto.Marshal(data)
	if err != nil {
		return nil, xerror.NewRuntimeError("error marshaling device data").Wrap(err)
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

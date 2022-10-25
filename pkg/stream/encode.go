package stream

import (
	"github.com/ntsd/cross-clipboard/pkg/device"
	"github.com/ntsd/cross-clipboard/pkg/protobuf"
	"github.com/ntsd/cross-clipboard/pkg/xerror"
	"google.golang.org/protobuf/proto"
)

// encodeClipboardData encode data for stream package `| data size (int 4 bytes) | data type (enum 1 byte) | protobuf message (struct n bytes) |`
func (s *StreamHandler) encodeClipboardData(dv *device.Device, clipboardData *protobuf.ClipboardData) ([]byte, error) {
	packageData := []byte{}

	// create proto clipboard data
	clipboardDataBytes, err := proto.Marshal(clipboardData)
	if err != nil {
		return nil, xerror.NewRuntimeError("error marshaling clipboard data").Wrap(err)
	}

	// encrypt clipboard data
	clipboardDataEncrypted, err := dv.PgpEncrypter.EncryptMessage(clipboardDataBytes)
	if err != nil {
		return nil, xerror.NewRuntimeError("error to encrypt clipboard data").Wrap(err)
	}
	encryptedDataSize := len(clipboardDataEncrypted)

	// append data size + 1 bytes for data type
	packageData = append(packageData, intToBytes(encryptedDataSize+1)...)
	// append data type
	packageData = append(packageData, byte(DataTypeClipboard))
	// append message
	packageData = append(packageData, clipboardDataEncrypted...)

	return packageData, nil
}

// encodeDeviceData encode data for stream package `| data size (int 4 bytes) | data type (enum 1 byte) | protobuf message (struct n bytes) |`
func (s *StreamHandler) encodeDeviceData(data *protobuf.DeviceData) ([]byte, error) {
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
	packageData = append(packageData, byte(DataTypeDevice))
	// append message
	packageData = append(packageData, deviceDataBytes...)

	return packageData, nil
}

// encodeSignal encode signal type `| data size (int 4 bytes) | signal type |`
func (s *StreamHandler) encodeSignal(signal Signal) ([]byte, error) {
	packageData := []byte{}

	// append data size + 1 bytes for data type
	packageData = append(packageData, intToBytes(1)...)
	// append DATA TYPE
	packageData = append(packageData, byte(signal))

	return packageData, nil
}

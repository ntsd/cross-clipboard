package clipboard

import (
	"time"

	"github.com/ntsd/cross-clipboard/pkg/device"
	"github.com/ntsd/cross-clipboard/pkg/protobuf"
)

// Clipboard struct for clipboard
type Clipboard struct {
	IsImage bool
	Data    []byte
	Size    uint32
	Time    time.Time
	Device  *device.Device
}

// ToProtobuf convert Clipboard to protocol buffer ClipboardData
func (c Clipboard) ToProtobuf() *protobuf.ClipboardData {
	return &protobuf.ClipboardData{
		IsImage:  c.IsImage,
		Data:     c.Data,
		DataSize: c.Size,
		Time:     c.Time.Unix(),
	}
}

// FromProtobuf convert protobuf.ClipboardData to Clipboard struct
func FromProtobuf(cd *protobuf.ClipboardData) Clipboard {
	return Clipboard{
		IsImage: cd.IsImage,
		Data:    cd.Data,
		Size:    cd.DataSize,
		Time:    time.UnixMicro(cd.Time),
	}
}

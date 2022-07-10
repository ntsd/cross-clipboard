package clipboard

import (
	"time"

	"github.com/ntsd/cross-clipboard/pkg/device"
)

// Clipboard struct for clipboard
type Clipboard struct {
	IsImage bool
	Data    []byte
	Size    uint32
	Time    time.Time
	Device  *device.Device
}

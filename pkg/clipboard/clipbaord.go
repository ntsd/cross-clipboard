package clipboard

import (
	"time"

	"github.com/ntsd/cross-clipboard/pkg/p2p"
)

// Clipboard struct for clipboard
type Clipboard struct {
	IsImage bool
	Data    []byte
	Size    uint32
	Time    time.Time
	Peer    *p2p.Peer
}

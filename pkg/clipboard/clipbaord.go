package clipboard

import (
	"time"

	"github.com/ntsd/cross-clipboard/pkg/p2p"
)

type Clipboard struct {
	Text []byte
	Size int
	Time time.Time
	Peer *p2p.Peer
}

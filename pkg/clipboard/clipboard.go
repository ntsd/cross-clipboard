package clipboard

import (
	"bytes"
	"context"
	"sync"

	"github.com/ntsd/cross-clipboard/pkg/config"
	"golang.design/x/clipboard"
)

type Clipboard struct {
	Config           config.Config
	ReadChannel      <-chan []byte
	Clipboards       [][]byte
	CurrentClipboard []byte
	mu               sync.RWMutex
}

func NewClipboard(cfg config.Config, clipboards [][]byte) *Clipboard {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	ch := clipboard.Watch(context.Background(), clipboard.FmtText)

	return &Clipboard{
		Config:      cfg,
		ReadChannel: ch,
		Clipboards:  clipboards,
	}
}

func limitAppendRotate[T any](limit int, slice []T, new T) []T {
	l := len(slice)
	if l >= limit {
		slice = slice[1:]
	}
	slice = append(slice, new)
	return slice
}

func (c *Clipboard) Write(newClipboard []byte) {
	if bytes.Compare(c.CurrentClipboard, newClipboard) != 0 {
		limitAppendRotate(c.Config.MaxHistory, c.Clipboards, newClipboard)
		clipboard.Write(clipboard.FmtText, newClipboard)
	}
}

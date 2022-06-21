package clipboard

import (
	"bytes"
	"context"
	"sync"

	"github.com/ntsd/cross-clipboard/pkg/config"
	"golang.design/x/clipboard"
)

type ClipboardManager struct {
	Config            config.Config
	ReadChannel       <-chan []byte
	Clipboards        [][]byte
	ClipboardsChannel chan [][]byte
	CurrentClipboard  []byte
	mu                sync.RWMutex
}

func NewClipboardManager(cfg config.Config) *ClipboardManager {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	ch := clipboard.Watch(context.Background(), clipboard.FmtText)

	return &ClipboardManager{
		Config:      cfg,
		ReadChannel: ch,
		Clipboards:  [][]byte{},
	}
}

func limitAppend[T any](limit int, slice []T, new T) []T {
	l := len(slice)
	if l >= limit {
		slice = slice[1:]
	}
	slice = append(slice, new)
	return slice
}

func (c *ClipboardManager) Write(newClipboard []byte) {
	if bytes.Compare(c.CurrentClipboard, newClipboard) != 0 {
		c.Clipboards = limitAppend(c.Config.MaxHistory, c.Clipboards, newClipboard)
		c.ClipboardsChannel <- c.Clipboards
		clipboard.Write(clipboard.FmtText, newClipboard)
	}
}

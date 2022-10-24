package clipboard

import (
	"bytes"
	"context"

	"github.com/ntsd/cross-clipboard/pkg/config"
	"golang.design/x/clipboard"
)

// ClipboardManager struct for clipbaord manager
type ClipboardManager struct {
	config *config.Config

	ReadTextChannel   <-chan []byte
	ReadImageChannel  <-chan []byte
	Clipboards        []Clipboard
	ClipboardsChannel chan []Clipboard
	CurrentClipboard  []byte
}

// NewClipboardManager create new clipbaord manager
func NewClipboardManager(cfg *config.Config) *ClipboardManager {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	textCh := clipboard.Watch(context.Background(), clipboard.FmtText)
	imgCh := clipboard.Watch(context.Background(), clipboard.FmtImage)

	return &ClipboardManager{
		config:            cfg,
		ReadTextChannel:   textCh,
		ReadImageChannel:  imgCh,
		ClipboardsChannel: make(chan []Clipboard),
		Clipboards:        []Clipboard{},
	}
}

// limitAppend append and rotate when limit
func limitAppend[T any](limit int, slice []T, new T) []T {
	l := len(slice)
	if l >= limit {
		slice = slice[1:]
	}
	slice = append(slice, new)
	return slice
}

// WriteClipboard write os clipbaord
func (c *ClipboardManager) WriteClipboard(newClipboard Clipboard) {
	if bytes.Compare(c.CurrentClipboard, newClipboard.Data) != 0 {
		// TODO avoid clipboard read channel after write by this
		if newClipboard.IsImage {
			clipboard.Write(clipboard.FmtImage, newClipboard.Data)
			return
		}
		clipboard.Write(clipboard.FmtText, newClipboard.Data)
	}
}

// AddClipboard add clipbaord to clipbaord history
func (c *ClipboardManager) AddClipboard(newClipboard Clipboard) {
	if bytes.Compare(c.CurrentClipboard, newClipboard.Data) != 0 {
		c.CurrentClipboard = newClipboard.Data
		c.Clipboards = limitAppend(c.config.MaxHistory, c.Clipboards, newClipboard)
		c.ClipboardsChannel <- c.Clipboards
	}
}

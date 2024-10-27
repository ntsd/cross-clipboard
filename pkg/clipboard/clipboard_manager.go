package clipboard

import (
	"bytes"
	"context"

	"github.com/ntsd/cross-clipboard/pkg/config"
	"github.com/ntsd/cross-clipboard/pkg/device"
	"golang.design/x/clipboard"
)

// ClipboardManager struct for clipbaord manager
type ClipboardManager struct {
	config *config.Config

	ReadTextChannel          <-chan []byte
	ReadImageChannel         <-chan []byte
	ClipboardsHistory        []*Clipboard
	ClipboardsHistoryUpdated chan struct{}
	receivedClipboard        *Clipboard
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
		config:                   cfg,
		ReadTextChannel:          textCh,
		ReadImageChannel:         imgCh,
		ClipboardsHistoryUpdated: make(chan struct{}),
		ClipboardsHistory:        []*Clipboard{},
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
	c.receivedClipboard = &newClipboard

	c.AddClipboardToHistory(&newClipboard)

	if newClipboard.IsImage {
		clipboard.Write(clipboard.FmtImage, newClipboard.Data)
		return
	}
	clipboard.Write(clipboard.FmtText, newClipboard.Data)
}

// AddClipboardToHistory add clipbaord to clipbaord history
func (c *ClipboardManager) AddClipboardToHistory(newClipboard *Clipboard) {
	c.ClipboardsHistory = limitAppend(c.config.MaxHistory, c.ClipboardsHistory, newClipboard)
	c.ClipboardsHistoryUpdated <- struct{}{}
}

// IsReceivedDevice returns true if it's the same device with the received clipboard
func (c *ClipboardManager) IsReceivedDevice(dv *device.Device) bool {
	if c.receivedClipboard == nil {
		return false
	}

	if c.receivedClipboard.Device == nil {
		return false
	}

	return c.receivedClipboard.Device.AddressInfo.ID == dv.AddressInfo.ID
}

// IsReceivedClipboard returns true if it's same clipboard data with the received clipboard
func (c *ClipboardManager) IsReceivedClipboard(clipboardData []byte) bool {
	if c.receivedClipboard == nil {
		return false
	}

	return bytes.Equal(clipboardData, c.receivedClipboard.Data)
}

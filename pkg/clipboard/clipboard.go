package clipboard

import (
	"context"
	"sync"

	"golang.design/x/clipboard"
)

type Clipboard struct {
	ReadChannel      <-chan []byte
	CurrentClipboard []byte
	mu               sync.RWMutex
}

func NewClipboard() *Clipboard {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	ch := clipboard.Watch(context.Background(), clipboard.FmtText)

	return &Clipboard{
		ReadChannel: ch,
	}
}

func (c *Clipboard) Write(newClipboard []byte) {
	c.mu.Lock()
	clipboard.Write(clipboard.FmtText, newClipboard)
	c.mu.Unlock()
}

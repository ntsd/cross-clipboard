package clipboard

import (
	"bytes"
	"context"

	"golang.design/x/clipboard"
)

type Clipboard struct {
	ReadChannel      <-chan []byte
	CurrentClipboard []byte
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
	if bytes.Compare(c.CurrentClipboard, newClipboard) != 0 {
		<-clipboard.Write(clipboard.FmtText, newClipboard)
	}
}

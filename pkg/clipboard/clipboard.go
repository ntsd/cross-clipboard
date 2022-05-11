package clipboard

import (
	"context"

	"golang.design/x/clipboard"
)

type Clipboard struct {
	ReadChannel <-chan []byte
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

func (c *Clipboard) Write(bytes []byte) {
	clipboard.Write(clipboard.FmtText, bytes)
}

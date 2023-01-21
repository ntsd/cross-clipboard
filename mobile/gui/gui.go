package gui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ntsd/cross-clipboard/pkg/crossclipboard"
)

type GUI struct {
	size image.Point

	cc crossclipboard.CrossClipboard
}

func NewGUI(cc crossclipboard.CrossClipboard) ebiten.Game {
	return &GUI{
		cc: cc,
	}
}

package gui

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

type GUI struct {
	size image.Point
}

func NewGUI() ebiten.Game {
	return &GUI{}
}

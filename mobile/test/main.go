package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ntsd/cross-clipboard/mobile/gui"
)

func main() {
	ebiten.SetWindowTitle("Cross Clipboard (Test)")
	ebiten.SetWindowSize(360, 800)

	g := gui.NewGUI()

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

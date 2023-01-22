package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ntsd/cross-clipboard/mobile/gui"
	"github.com/ntsd/cross-clipboard/pkg/config"
	"github.com/ntsd/cross-clipboard/pkg/crossclipboard"
)

func main() {
	ebiten.SetWindowTitle("Cross Clipboard (Test)")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	cc, err := crossclipboard.NewCrossClipboard(cfg)
	if err != nil || cc == nil {
		log.Fatal(err)
	}

	g := gui.NewGUI(*cc)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

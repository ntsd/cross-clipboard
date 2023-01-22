package mobile

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/mobile"
	"github.com/ntsd/cross-clipboard/mobile/gui"
	"github.com/ntsd/cross-clipboard/pkg/config"
	"github.com/ntsd/cross-clipboard/pkg/crossclipboard"
)

func init() {
	ebiten.SetWindowTitle("Cross Clipboard (Beta)")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	cc, err := crossclipboard.NewCrossClipboard(cfg)
	if err != nil || cc == nil {
		log.Fatal(err)
	}

	g := gui.NewGUI(*cc)

	mobile.SetGame(g)
}

// Dummy is a dummy exported function.
//
// gomobile doesn't compile a package that doesn't include any exported function.
// Dummy forces gomobile to compile this package.
func Dummy() {}

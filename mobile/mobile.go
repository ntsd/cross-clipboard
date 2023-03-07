package mobile

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/mobile"
	"github.com/ntsd/cross-clipboard/mobile/gui"
)

func init() {
	ebiten.SetWindowTitle("Cross Clipboard (Test)")

	g := gui.NewGUI()

	mobile.SetGame(g)
}

// Dummy is a dummy exported function.
//
// gomobile doesn't compile a package that doesn't include any exported function.
// Dummy forces gomobile to compile this package.
func Dummy() {}

package gui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Draw renders the terminal GUI to the ebtien window. Required to implement the ebiten interface.
func (g *GUI) Draw(screen *ebiten.Image) {
	// Blink the cursor.
	t := g.text
	if g.counter%60 < 30 {
		t += "_"
	}
	ebitenutil.DebugPrint(screen, t)
}

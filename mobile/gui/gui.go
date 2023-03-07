package gui

import (
	"image"

	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/ntsd/cross-clipboard/pkg/config"
	"github.com/ntsd/cross-clipboard/pkg/crossclipboard"
)

type GUI struct {
	ui *ebitenui.UI

	// gui data
	size       image.Point
	startError error

	// cross clipboard
	cc crossclipboard.CrossClipboard
}

func NewGUI() ebiten.Game {
	gui := &GUI{}

	ui, err := createUI()
	if err != nil {
		gui.startError = err
		return gui
	}
	gui.ui = ui

	// TODO fix load config
	cfg, err := config.LoadConfig()
	if err != nil {
		gui.startError = err
		return gui
	}

	cc, err := crossclipboard.NewCrossClipboard(cfg)
	if err != nil || cc == nil {
		gui.startError = err
		return gui
	}

	gui.cc = *cc

	return gui
}

// Update updates a game by one tick.
func (g *GUI) Update() error {
	g.ui.Update()
	return nil
}

func (g GUI) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	w, h := outsideWidth, outsideHeight

	if g.size.X != w || g.size.Y != h {
		g.size = image.Point{X: w, Y: h}
	}

	return w, h
}

// Draw draw the game screen. The given argument represents a screen image.
func (g *GUI) Draw(screen *ebiten.Image) {
	// print debug eror when starting errors
	// if g.startError != nil {
	// 	ebitenutil.DebugPrint(screen, g.startError.Error())
	// 	return
	// }

	g.ui.Draw(screen)
}

package gui

import "image"

func (g GUI) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	w, h := outsideWidth, outsideHeight

	if g.size.X != w || g.size.Y != h {
		g.size = image.Point{X: w, Y: h}
		g.resize(w, h)
	}

	return w, h
}

func (g *GUI) resize(w, h int) {

}

package gui

import (
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *GUI) Update() error {
	// Add runes that are input by the user by AppendInputChars.
	// Note that AppendInputChars result changes every frame, so you need to call this
	// every frame.
	g.runes = ebiten.AppendInputChars(g.runes[:0])
	g.text += string(g.runes)

	// Adjust the string to be at most 10 lines.
	ss := strings.Split(g.text, "\n")
	if len(ss) > 10 {
		g.text = strings.Join(ss[len(ss)-10:], "\n")
	}

	// If the enter key is pressed, add a line break.
	if repeatingKeyPressed(ebiten.KeyEnter) || repeatingKeyPressed(ebiten.KeyNumpadEnter) {
		g.text += "\n"
	}

	// If the backspace key is pressed, remove one character.
	if repeatingKeyPressed(ebiten.KeyBackspace) {
		if len(g.text) >= 1 {
			g.text = g.text[:len(g.text)-1]
		}
	}

	g.counter++
	return nil
}

// repeatingKeyPressed return true when key is pressed considering the repeat state.
func repeatingKeyPressed(key ebiten.Key) bool {
	const (
		delay    = 30
		interval = 3
	)
	d := inpututil.KeyPressDuration(key)
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}

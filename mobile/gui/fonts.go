package gui

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

type Fonts struct {
	face         font.Face
	titleFace    font.Face
	bigTitleFace font.Face
	toolTipFace  font.Face
}

func loadFonts() (*Fonts, error) {
	fontFace, err := loadFont(20)
	if err != nil {
		return nil, err
	}

	titleFontFace, err := loadFont(24)
	if err != nil {
		return nil, err
	}

	bigTitleFontFace, err := loadFont(28)
	if err != nil {
		return nil, err
	}

	toolTipFace, err := loadFont(15)
	if err != nil {
		return nil, err
	}

	return &Fonts{
		face:         fontFace,
		titleFace:    titleFontFace,
		bigTitleFace: bigTitleFontFace,
		toolTipFace:  toolTipFace,
	}, nil
}

func loadFont(size float64) (font.Face, error) {
	ttfFont, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(ttfFont, &truetype.Options{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	}), nil
}

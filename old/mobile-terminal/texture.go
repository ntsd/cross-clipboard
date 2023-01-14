//go:build darwin || linux || windows
// +build darwin linux windows

package main

import (
	"fmt"
	"log"

	"image"
	_ "image/png" // The _ means to import a package purely for its initialization side effects.

	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/gl"
)

var nilTexture = gl.Texture{Value: 0xFFFFFFFF}

func loadTexture(glc gl.Context, name string, yflip bool) (gl.Texture, *image.NRGBA, error) {

	imgFile := name
	imgIn, errImg := asset.Open(imgFile)
	if errImg != nil {
		return nilTexture, nil, fmt.Errorf("open texture image: %s: %v", imgFile, errImg)
	}
	img, _, errDec := image.Decode(imgIn)
	if errDec != nil {
		return nilTexture, nil, fmt.Errorf("decode texture image: %s: %v", imgFile, errDec)
	}
	if img == nil {
		return nilTexture, nil, fmt.Errorf("decode texture image: %s: nil", imgFile)
	}
	log.Printf("texture image loaded: %s", imgFile)
	i, ok := img.(*image.NRGBA)
	if !ok {
		return nilTexture, nil, fmt.Errorf("unexpected image type: %s: %v", imgFile, img.ColorModel())
	}

	t, errUpload := uploadImage(glc, name, i, yflip)

	return t, i, errUpload
}

func uploadImage(glc gl.Context, name string, i *image.NRGBA, yflip bool) (gl.Texture, error) {
	b := i.Bounds()
	log.Printf("NRGBA image: %s %dx%d", name, b.Max.X, b.Max.Y)
	if yflip {
		flipY(name, i)
	}

	t := glc.CreateTexture()
	glc.BindTexture(gl.TEXTURE_2D, t)
	glc.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	glc.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	glc.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	glc.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	bounds := i.Bounds()
	w := bounds.Max.X - bounds.Min.X
	h := bounds.Max.Y - bounds.Min.Y

	// https://godoc.org/golang.org/x/mobile/gl
	// TexImage2D(target Enum, level int, internalFormat int, width, height int, format Enum, ty Enum, data []byte)
	glc.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, w, h, gl.RGBA, gl.UNSIGNED_BYTE, i.Pix)

	log.Printf("texture image uploaded: %s %dx%d", name, w, h)

	return t, nil
}

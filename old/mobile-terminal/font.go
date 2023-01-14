//go:build darwin || linux || windows
// +build darwin linux windows

package main

import (
	//"fmt"
	"bytes"
	"log"

	//"os"
	"encoding/binary"
	"image"
	"image/color"

	//"image/png"

	//"golang.org/x/mobile/asset"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/gl"

	"github.com/udhos/pixfont"
)

type fontAtlas struct {
	glc gl.Context
	tex gl.Texture
	//vert      gl.Buffer
	//elem      gl.Buffer
	//elemCount int
	coordVer gl.Attrib
	coordTex gl.Attrib
}

type fontText struct {
	vert      gl.Buffer
	elem      gl.Buffer
	elemCount int
	atlas     *fontAtlas
}

func newText(atlas *fontAtlas) *fontText {
	glc := atlas.glc
	t := &fontText{}
	t.atlas = atlas
	t.vert = glc.CreateBuffer()
	t.elem = glc.CreateBuffer()
	return t
}

func (t *fontText) delete() {
	t.atlas.glc.DeleteBuffer(t.vert)
	t.atlas.glc.DeleteBuffer(t.elem)
}

func (t *fontText) write(s string) {

	glc := t.atlas.glc // shortcut

	v := make([]float32, 0, 4*5*len(s))
	e := make([]uint32, 0, len(s))

	for i, b := range s {

		x1 := float32(i)
		x2 := x1 + 1.0

		c := b - fontFirst
		unit := 1.0 / float32(fontCount)
		s1 := float32(c) * unit
		s2 := float32(c+1) * unit

		v = append(v,
			// vert       texture
			// ---------  -------
			x1, 1.0, 0.0, s1, 1.0, // 0
			x1, 0.0, 0.0, s1, 0.0, // 1
			x2, 0.0, 0.0, s2, 0.0, // 2
			x2, 1.0, 0.0, s2, 1.0, // 3
		)

		j := 4 * uint32(i)

		e = append(e,
			j, j+1, j+2, // triangle 1
			j+2, j+3, j, // triangle 2
		)

	}
	t.elemCount = len(e)

	bytesV := f32.Bytes(binary.LittleEndian, v...)
	bytesE := intsToBytes(e)

	glc.BindBuffer(gl.ARRAY_BUFFER, t.vert)
	glc.BufferData(gl.ARRAY_BUFFER, bytesV, gl.DYNAMIC_DRAW)

	glc.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, t.elem)
	glc.BufferData(gl.ELEMENT_ARRAY_BUFFER, bytesE, gl.DYNAMIC_DRAW)
}

func (t *fontText) draw() {

	glc := t.atlas.glc // shortcut

	glc.BindBuffer(gl.ARRAY_BUFFER, t.vert)
	glc.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, t.elem)

	elemFirst := 0
	elemCount := t.elemCount
	elemType := gl.Enum(gl.UNSIGNED_INT)
	elemSize := 4 // 4-byte int

	strideSize := 5 * 4 // 5 x 4 bytes (5 x 4-byte float)
	itemsPosition := 3
	itemsTexture := 2
	offsetPosition := 0
	offsetTexture := itemsPosition * 4 // 3 x 4 bytes

	glc.VertexAttribPointer(t.atlas.coordVer, itemsPosition, gl.FLOAT, false, strideSize, offsetPosition)
	glc.VertexAttribPointer(t.atlas.coordTex, itemsTexture, gl.FLOAT, false, strideSize, offsetTexture)

	glc.BindTexture(gl.TEXTURE_2D, t.atlas.tex)

	glc.DrawElements(gl.TRIANGLES, elemCount, elemType, elemFirst*elemSize)
}

const (
	fontFirst   = 32
	fontPastend = 127
	fontCount   = fontPastend - fontFirst
)

func newAtlas(glc gl.Context, c color.Color, coordVert, coordTex gl.Attrib) (*fontAtlas, error) {

	first := fontFirst
	pastend := fontPastend
	size := fontCount
	b := make([]byte, 0, size)
	buf := bytes.NewBuffer(b)

	for i := first; i < pastend; i++ {
		err := buf.WriteByte(byte(i))
		if err != nil {
			log.Printf("newAtlas: %d %v", i, err)
		}
	}

	str := buf.String()

	log.Printf("newAtlas: chars=%d str=%d: [%v]", size, len(str), str)

	pixfont.Spacing = 0 // default: 1 pixel

	width := pixfont.MeasureString(str)
	log.Printf("newAtlas: width=%d", width)

	img := image.NewNRGBA(image.Rect(0, 0, width, 8))

	bo := img.Bounds()
	fontWidth := width / size
	fontHeight := bo.Max.Y
	log.Printf("newAtlas: atlas=%dx%d => font=%dx%d", bo.Max.X, bo.Max.Y, fontWidth, fontHeight)

	//c := color.NRGBA{128,230,128,255}
	w := pixfont.DrawString(img, 0, 0, str, c)

	log.Printf("newAtlas: drawn %d pixels", w)

	/*
		f, errWr := os.OpenFile("atlas.png", os.O_CREATE|os.O_RDWR, 0644)
		if errWr != nil {
			log.Printf("IO: %v", errWr)
		}
		png.Encode(f, img)
		//f.Flush()
		f.Close()
	*/

	tex, errUpload := uploadImage(glc, "<fontAtlas>", img, true)
	if errUpload != nil {
		return nil, errUpload
	}

	a := &fontAtlas{}

	a.glc = glc
	a.tex = tex
	//a.vert = glc.CreateBuffer()
	//a.elem = glc.CreateBuffer()
	a.coordVer = coordVert
	a.coordTex = coordTex

	return a, nil
}

func (a *fontAtlas) delete() {
	a.glc.DeleteTexture(a.tex)
	//a.glc.DeleteBuffer(a.vert)
	//a.glc.DeleteBuffer(a.elem)
}

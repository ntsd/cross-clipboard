//go:build darwin || linux || windows
// +build darwin linux windows

package main

import (
	"image/color"
	"log"

	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/gl"
)

type terminal struct {
	width      int
	height     int
	gl         gl.Context
	programTex gl.Program

	// texturizing shader
	texPosition     gl.Attrib
	texTextureCoord gl.Attrib

	atlas    *fontAtlas
	textFont *fontText

	shaderTexVert string
	shaderTexFrag string
}

func newTerminal() (*terminal, error) {
	t := &terminal{}

	if errVert := flagStr(&t.shaderTexVert, "shader_tex.vert"); errVert != nil {
		log.Printf("load vertex tex shader: %v", errVert)
		return nil, errVert
	}

	if errFrag := flagStr(&t.shaderTexFrag, "shader_tex.frag"); errFrag != nil {
		log.Printf("load fragment tex shader: %v", errFrag)
		return nil, errFrag
	}

	return t, nil
}

func (t *terminal) start(glc gl.Context) {
	// create program
	programTex, err := glutil.CreateProgram(glc, t.shaderTexVert, t.shaderTexFrag)
	if err != nil {
		log.Printf("start: error creating GL texturizer program: %v", err)
		return
	}
	t.programTex = programTex
	log.Printf("start: texturizing shader compiled")

	// create text position and texture cord
	t.texPosition = getAttribLocation(glc, t.programTex, "position")
	t.texTextureCoord = getAttribLocation(glc, t.programTex, "textureCoord")

	// create text font
	t.atlas, err = newAtlas(glc, color.NRGBA{128, 230, 128, 255}, t.texPosition, t.texTextureCoord)
	if err != nil {
		log.Printf("new font error: %v", err)
		return
	}

	t.textFont = newText(t.atlas)
	t.textFont.write("Test text")

	glc.ClearColor(.5, .5, .5, 1) // gray background
	glc.ClearDepthf(1)            // default
	glc.Enable(gl.DEPTH_TEST)     // enable depth testing
	glc.DepthFunc(gl.LEQUAL)      // gl.LESS is default depth test
	glc.DepthRangef(0, 1)         // default

	t.gl = glc

	log.Printf("start: shaders initialized")
}

func (t *terminal) stop() {
	log.Panicln("stop terminal")
}

func (t *terminal) paint() {
	// clear screen
	t.gl.ClearColor(1, 1, 1, 1)
	t.gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// enable blend
	t.gl.Enable(gl.BLEND)
	t.gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	// prepare text program
	t.gl.UseProgram(t.programTex)
	t.gl.EnableVertexAttribArray(t.texPosition)
	t.gl.EnableVertexAttribArray(t.texTextureCoord)

	// draw text font
	t.textFont.draw()

	// clean-up
	t.gl.DisableVertexAttribArray(t.texPosition)
	t.gl.DisableVertexAttribArray(t.texTextureCoord)

	t.gl.Disable(gl.BLEND)
}

func getAttribLocation(glc gl.Context, prog gl.Program, attr string) gl.Attrib {
	location := glc.GetAttribLocation(prog, attr)
	// FIXME 1000 is a hack to detect a bad location.Value, since it can't represent -1
	if location.Value > 1000 {
		log.Printf("bad attribute '%s' location: %d", attr, location.Value)
	}
	return location
}

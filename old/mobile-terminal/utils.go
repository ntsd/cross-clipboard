//go:build darwin || linux || windows
// +build darwin linux windows

package main

import (
	"bytes"
	"encoding/binary"
	"image"
	"io/ioutil"
	"log"
	"strings"

	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/gl"
)

func intsToBytes(s []uint32) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, s)
	b := buf.Bytes()
	return b
}

func flipY(name string, img *image.NRGBA) {
	b := img.Bounds()
	midY := (b.Max.Y - b.Min.Y) / 2
	for x := b.Min.X; x < b.Max.X; x++ {
		for y1 := b.Min.Y; y1 < midY; y1++ {
			y2 := b.Max.Y - y1 - 1
			c1 := img.At(x, y1)
			c2 := img.At(x, y2)
			img.Set(x, y1, c2)
			img.Set(x, y2, c1)
		}
	}
	log.Printf("image y-flipped: %s", name)
}

func flagBool(value *bool, name string) {
	*value = exists(name)
	log.Printf("flagBool: %s = %v", name, *value)
}

func flagStr(value *string, name string) error {
	b, errLoad := loadFull(name)
	if errLoad != nil {
		log.Printf("flagStr: %s: %v", name, errLoad)
		return errLoad
	}
	*value = strings.TrimSpace(string(b))
	log.Printf("flagStr: %s = [%v]", name, *value)
	return nil
}

func exists(name string) bool {
	f, errOpen := asset.Open(name)
	if errOpen != nil {
		return false
	}
	f.Close()
	return true
}

func loadFull(name string) ([]byte, error) {
	f, errOpen := asset.Open(name)
	if errOpen != nil {
		return nil, errOpen
	}
	defer f.Close()
	buf, errRead := ioutil.ReadAll(f)
	if errRead != nil {
		return nil, errRead
	}
	log.Printf("loaded: %s (%d bytes)", name, len(buf))
	return buf, nil
}

func getUniformLocation(glc gl.Context, prog gl.Program, uniform string) gl.Uniform {
	location := glc.GetUniformLocation(prog, uniform)
	if location.Value < 0 {
		log.Printf("bad uniform '%s' location: %d", uniform, location.Value)
	}
	return location
}

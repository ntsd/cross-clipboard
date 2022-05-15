//go:build darwin || linux || windows
// +build darwin linux windows

// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// An app that draws a green triangle on a red background.
//
// In order to build this program as an Android APK, using the gomobile tool.
//
// See http://godoc.org/golang.org/x/mobile/cmd/gomobile to install gomobile.
//
// Get the basic example and use gomobile to build or install it on your device.
//
//	$ go get -d golang.org/x/mobile/example/basic
//	$ gomobile build golang.org/x/mobile/example/basic # will build an APK
//
//	# plug your Android device to your computer or start an Android emulator.
//	# if you have adb installed on your machine, use gomobile install to
//	# build and deploy the APK to an Android target.
//	$ gomobile install golang.org/x/mobile/example/basic
//
// Switch to your device or emulator to start the Basic application from
// the launcher.
// You can also run the application on your desktop by running the command
// below. (Note: It currently doesn't work on Windows.)
//
//	$ go install golang.org/x/mobile/example/basic && basic
package main

import (
	"github.com/ntsd/cross-clipboard/pkg/p2p"
	"github.com/ntsd/cross-clipboard/pkg/utils"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/app/debug"
	"golang.org/x/mobile/gl"
)

var (
	fps *debug.FPS

	touchX float32
	touchY float32
)

func main() {
	app.Main(func(a app.App) {
		var glctx gl.Context

		cfg := utils.Config{
			RendezvousString: "default-group",
			ProtocolID:       "/cross-clipboard/0.0.1",
			ListenHost:       "0.0.0.0",
			ListenPort:       4001,
		}
		p2p.StartP2P(cfg)

		var sz size.Event
		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					glctx, _ = e.DrawContext.(gl.Context)
					onStart(glctx)
					a.Send(paint.Event{}) // keep animating
				case lifecycle.CrossOff:
					onStop(glctx)
					glctx = nil
				}
			case paint.Event:
				if glctx == nil || e.External {
					continue
				}
				onPaint(glctx, sz)
				a.Publish()
				a.Send(paint.Event{}) // keep animating
			case size.Event:
				sz = e
				touchX = float32(sz.WidthPx / 2)
				touchY = float32(sz.HeightPx / 2)
			case touch.Event:
				touchX = e.X
				touchY = e.Y
			}
		}
	})
}

func onStart(glctx gl.Context) {
}
func onStop(glctx gl.Context) {
	fps.Release()
}
func onPaint(glctx gl.Context, sz size.Event) {
	glctx.ClearColor(1, 0, 0, 1)
	glctx.Clear(gl.COLOR_BUFFER_BIT)
}

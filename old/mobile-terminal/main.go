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
	"log"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/gl"
)

func main() {
	app.Main(func(a app.App) {
		t, err := newTerminal()
		if err != nil {
			log.Fatalf("error creating terminal: %v", err)
		}

	LOOP:
		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					glc, isGL := e.DrawContext.(gl.Context)
					if !isGL {
						log.Printf("Lifecycle: visible: bad GL context")
						continue LOOP
					}

					t.start(glc)

					a.Send(paint.Event{}) // keep animating
				case lifecycle.CrossOff:
					t.stop()
				}
			case paint.Event:
				if t.gl == nil || e.External {
					continue
				}

				t.paint()

				a.Publish()
				a.Send(paint.Event{}) // keep animating
			case size.Event:
				// listen screen size for screen rotate
			case touch.Event:
			}
		}
	})
}

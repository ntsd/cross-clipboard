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
package mobile

import (
	"github.com/ntsd/cross-clipboard/pkg/config"
	"github.com/ntsd/cross-clipboard/pkg/crossclipboard"
)

type CrossClipbardMobile struct {
}

func (c *CrossClipbardMobile) Start() {
	cfg := config.Config{
		GroupName:  "default",
		ProtocolID: "/cross-clipboard/0.0.1",
		ListenHost: "0.0.0.0",
		ListenPort: 4001,
	}
	crossclipboard.NewCrossClipboard(cfg)
}

func NewCrossClipbardMobile() *CrossClipbardMobile {
	return &CrossClipbardMobile{}
}

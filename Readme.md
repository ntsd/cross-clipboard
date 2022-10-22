# Cross Clipboard

A multi device clipboard sharing over P2P network.

![Cross Clipboard Preview](/docs/preview-home.jpg)

![Cross Clipboard Preview](/docs/preview-setting.jpg)

## Features

- Share text/image clipboard data (Done) - user can share clipboard data with other devices
- P2P connection (Done) - the device will connect to others using P2P connection
- Multicast DNS (Done) - discover a device in the same network
- E2E encryption (Done) - encrypt the clipboard data using OpenPGP
- Cross-platform desktop (Done) - support Windows, Linux and Darwin (macOS)
- Terminal GUI (Ongoing) - terminal user interface for end user
- Cross-platform mobile (Plan) - support iOS and Android

## Technology

- Go
- libp2p
- Multicast DNS (mDNS)
- Protobuf

## Installation

### Go install

for Go user you can just install using go package

```shell
go install github.com/ntsd/cross-clipboard@latest
```

### Headless Linux

for headless linux you might need to install `xvfb`.

```shell
# install libx11-dev abd Xvfb
sudo apt install -y libx11-dev xvfb

# initialize a virtual frame buffer (can put in .profile)
Xvfb :99 -screen 0 1024x768x24 > /dev/null 2>&1 &
export DISPLAY=:99.0
```

## Run

`go run cmd/cross-clipboard.go`

## Build

### Build Desktop

`go build ./cmd/cross-clipboard.go`

### Build Mobile (Plan)

- Install NDK >=21.3.6528147

- Install Go mobile

```shell
go install golang.org/x/mobile/cmd/gomobile@latest
```

`gomobile build ./mobile/...`

### Protobuf gen

Generate a protobuf go file using protoc

`protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative data.proto`

## TODO

- Release binary file (Ongoing)
- Publish to Homebrew and Deb package (Ongoing)
- Handle disconnect and error (Ongoing)
- Fix avoid sending data back to the peer it received from (Plan)
- Unit testing (Plan)

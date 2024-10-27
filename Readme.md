# Cross Clipboard

A multi device clipboard sharing over P2P network.

![Cross Clipboard Preview](/docs/preview-home.jpg)

![Cross Clipboard Preview](/docs/preview-setting.jpg)

## Features

- Share text/image clipboard data (Done) - user can share clipboard data with other devices.
- P2P connection (Done) - the device will connect to others using the P2P connection.
- Multicast DNS (Done) - discover a device in the same network with multicast DNS.
- E2E encryption (Done) - encrypt the clipboard data using OpenPGP.
- Terminal GUI (Done) - terminal user interface for the end user.
- Cross-platform desktop (Done) - support Windows, Linux, and Darwin (macOS).
- Cross-platform mobile (Ongoing) - support iOS and Android.

## Libraries 

- [libp2p](https://github.com/libp2p/go-libp2p)
- [clipboard](https://github.com/golang-design/clipboard)
- [tview](https://github.com/rivo/tview)
- [protobuf](https://developers.google.com/protocol-buffers)

## Installation

### Go install

for Go user you can just install using go package

Required: Go 1.23 or above

```shell
go install github.com/ntsd/cross-clipboard@latest
```

### Homebrew tap install

```shell
brew tap ntsd/cross-clipboard
brew install cross-clipboard
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

UI mode

`cross-clipboard`

Terminal mode

`cross-clipboard -t`

## Development

```shell
git clone https://github.com/ntsd/cross-clipboard
go run main.go
```

## Build

### Build Desktop

`go build .`

### Build Mobile (Ongoing)

- Install NDK >=21.3.6528147

- Install Go mobile

```shell
go install golang.org/x/mobile/cmd/gomobile@latest
```

`gomobile build mobile/...`

### Protobuf gen

Generate a protobuf go file using protoc

`protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative data.proto`

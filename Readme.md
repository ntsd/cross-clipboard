# Cross Clipboard

A multi device clipboard sharing over p2p network.

## TODO

- E2E Encryption (Done)
- Image clipboard (Done)
- Trust device (Done)
- Terminal GUI (Ongoing)
- Save device info to storage (Done)
- Command Line Tools
- Release binary file
- Publish to Homebrew and Deb package
- Handle disconnect and error
- Fix avoid sending data back to the peer it received from
- Mobile Support
- Unit tests

## Technology

- Go
- libp2p
- Multicast DNS (mDNS)
- Protobuf

## Installation

### Headless Linux

```shell
# install libx11-dev abd Xvfb
sudo apt install -y libx11-dev xvfb

# initialize a virtual frame buffer (can put in .profile)
Xvfb :99 -screen 0 1024x768x24 > /dev/null 2>&1 &
export DISPLAY=:99.0
```

### Mobile Build

- Install NDK >=21.3.6528147

- Install Go mobile

```shell
go install golang.org/x/mobile/cmd/gomobile@latest
```

## Run

`go run cmd/main.go`

## Proto gen

`protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative data.proto`

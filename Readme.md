# Cross Clipboard

A multi device clipboard sharing over p2p in lan network.

## TODO

- E2E Encryption (Done)
- Image clipboard (Done)
- Terminal GUI
- Avoid sending data back to the peer it received from
- Trust device
- Mobile Support
- Add private key pass phrase

## Technology

- Go
- libp2p
- Multicast DNS (mDNS)
- [Gogo Protobuf](https://github.com/gogo/protobuf)

## Installation

### Linux

```shell
# install libx11-dev
sudo apt install libx11-dev

# install Xvfb
sudo apt install -y xvfb

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

## Proto gen

`protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative data.proto`

# Cross Clipboard

A multi device clipboard sharing over p2p in lan network.

## TODO

- Terminal GUI
- Verify Peer
- E2E Encryption
- Mobile Support
- Image clipboard

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

## References

<https://github.com/rivo/tview>

<https://inblockchainwetrust.medium.com/running-go-from-android-ios-tutorial-7f1d456c5b0f>

<https://github.com/golang/go/wiki/Mobile>

<https://pkg.go.dev/golang.org/x/mobile/example>

<https://github.com/gomatcha/matcha>

<https://stackoverflow.com/questions/41607634/gomobile-callback-forward-realtime-downloaded-content-to-android>

<https://www.sajalkayan.com/post/android-apps-golang.html>

<https://sites.google.com/a/athaydes.com/renato-athaydes/posts/buildingamobilefrontendforagoapplicationusingflutter>

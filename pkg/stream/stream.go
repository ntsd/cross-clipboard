package stream

import (
	"bufio"
	"fmt"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/ntsd/cross-clipboard/pkg/clipboard"
	"github.com/ntsd/cross-clipboard/pkg/config"
	"github.com/ntsd/cross-clipboard/pkg/crypto"
	"github.com/ntsd/cross-clipboard/pkg/device"
	"github.com/ntsd/cross-clipboard/pkg/devicemanager"
)

// StreamHandler struct for stream handler
type StreamHandler struct {
	Config           config.Config
	ClipboardManager *clipboard.ClipboardManager
	DeviceManager    *devicemanager.DeviceManager
	LogChan          chan string
	ErrorChan        chan error

	pgpDecrypter *crypto.PGPDecrypter
}

// NewStreamHandler initial new stream handler
func NewStreamHandler(
	cfg config.Config,
	cp *clipboard.ClipboardManager,
	deviceManager *devicemanager.DeviceManager,
	logChan chan string,
	errorChan chan error,
	pgpDecrypter *crypto.PGPDecrypter,
) *StreamHandler {
	s := &StreamHandler{
		Config:           cfg,
		ClipboardManager: cp,
		DeviceManager:    deviceManager,
		LogChan:          logChan,
		ErrorChan:        errorChan,
		pgpDecrypter:     pgpDecrypter,
	}
	go s.CreateWriteData()
	return s
}

// HandleStream handler when a peer connect this host
func (s *StreamHandler) HandleStream(stream network.Stream) {
	s.LogChan <- fmt.Sprintf("peer %s connecting to this host", stream.Conn().RemotePeer())

	// Create a new peer
	dv := device.NewDevice(peer.AddrInfo{
		ID:    stream.Conn().RemotePeer(),
		Addrs: []multiaddr.Multiaddr{stream.Conn().RemoteMultiaddr()},
	}, stream)
	s.DeviceManager.AddDevice(dv)
	dv.Reader = bufio.NewReader(stream)
	dv.Writer = bufio.NewWriter(stream)

	go s.CreateReadData(dv.Reader, dv)

	s.LogChan <- fmt.Sprintf("peer %s connected to this host", stream.Conn().RemotePeer())
	// 'stream' will stay open until you close it (or the other side closes it).
}

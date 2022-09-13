package crossclipboard

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/multiformats/go-multiaddr"
	"github.com/ntsd/cross-clipboard/pkg/clipboard"
	"github.com/ntsd/cross-clipboard/pkg/config"
	"github.com/ntsd/cross-clipboard/pkg/crypto"
	"github.com/ntsd/cross-clipboard/pkg/device"
	"github.com/ntsd/cross-clipboard/pkg/devicemanager"
	"github.com/ntsd/cross-clipboard/pkg/discovery"
	"github.com/ntsd/cross-clipboard/pkg/stream"
	"github.com/ntsd/cross-clipboard/pkg/xerror"
)

// CrossClipboard cross clipbaord struct
type CrossClipboard struct {
	Host   host.Host
	Config config.Config

	ClipboardManager *clipboard.ClipboardManager
	DeviceManager    *devicemanager.DeviceManager

	LogChan chan string
	ErrChan chan error
}

// NewCrossClipboard initial cross clipbaord
func NewCrossClipboard(cfg config.Config) (*CrossClipboard, error) {
	cc := &CrossClipboard{
		Config:  cfg,
		LogChan: make(chan string),
		ErrChan: make(chan error),
	}

	cc.ClipboardManager = clipboard.NewClipboardManager(cc.Config)
	cc.DeviceManager = devicemanager.NewDeviceManager()

	ctx := context.Background()

	// 0.0.0.0 will listen on any interface device.
	sourceMultiAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", cc.Config.ListenHost, cc.Config.ListenPort))
	if err != nil {
		return nil, xerror.NewFatalError(err)
	}

	// libp2p.New constructs a new libp2p Host.
	host, err := libp2p.New(
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(cc.Config.ID),
	)
	if err != nil {
		return nil, xerror.NewFatalError(err)
	}
	cc.Host = host

	pgpDecrypter, err := crypto.NewPGPDecrypter(cfg.PGPPrivateKey)
	if err != nil {
		return nil, xerror.NewFatalError(err)
	}

	go func() {
		streamHandler := stream.NewStreamHandler(
			cc.Config,
			cc.ClipboardManager,
			cc.DeviceManager,
			cc.LogChan,
			cc.ErrChan,
			pgpDecrypter,
		)

		// Set a function as stream handler.
		// This function is called when a peer initiates a connection and starts a stream with this peer.
		cc.Host.SetStreamHandler(protocol.ID(cc.Config.ProtocolID), streamHandler.HandleStream)
		cc.LogChan <- fmt.Sprintf("[*] your multiaddress is: /ip4/%s/tcp/%v/p2p/%s", cc.Config.ListenHost, cc.Config.ListenPort, host.ID().Pretty())

		peerInfoChan, err := discovery.InitMultiMDNS(cc.Host, cc.Config.GroupName, cc.LogChan)
		if err != nil {
			cc.ErrChan <- xerror.NewFatalError(err)
		}

		for peerInfo := range peerInfoChan { // when discover a peer
			cc.LogChan <- fmt.Sprintf("connecting to peer host: %s", peerInfo)

			if err := cc.Host.Connect(ctx, peerInfo); err != nil {
				cc.ErrChan <- fmt.Errorf("connect error: %w", err)
				continue
			}

			// open a stream, this stream will be handled by handleStream other end
			stream, err := cc.Host.NewStream(ctx, peerInfo.ID, protocol.ID(cc.Config.ProtocolID))
			if err != nil {
				cc.ErrChan <- fmt.Errorf("new stream error: %w", err)
				continue
			}

			dv := device.NewDevice(peerInfo, stream)
			cc.DeviceManager.AddDevice(dv)
			go streamHandler.CreateReadData(dv.Reader, dv.AddressInfo.ID.Pretty())

			cc.LogChan <- fmt.Sprintf("connected to peer host: %s", peerInfo)
		}
	}()

	return cc, nil
}

package crossclipboard

import (
	"bufio"
	"context"
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
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
	Config *config.Config

	ClipboardManager *clipboard.ClipboardManager
	DeviceManager    *devicemanager.DeviceManager

	streamHandler *stream.StreamHandler

	LogChan   chan string
	ErrorChan chan error
}

// NewCrossClipboard initial cross clipbaord
func NewCrossClipboard(cfg *config.Config) (*CrossClipboard, error) {
	cc := &CrossClipboard{
		Config:    cfg,
		LogChan:   make(chan string),
		ErrorChan: make(chan error),
	}

	cc.ClipboardManager = clipboard.NewClipboardManager(cc.Config)
	cc.DeviceManager = devicemanager.NewDeviceManager(cc.Config)

	ctx := context.Background()

	// 0.0.0.0 will listen on any interface device.
	sourceMultiAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", cc.Config.ListenHost, cc.Config.ListenPort))
	if err != nil {
		return nil, xerror.NewFatalError("error to multiaddr.NewMultiaddr").Wrap(err)
	}

	// libp2p.New constructs a new libp2p Host.
	host, err := libp2p.New(
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(cc.Config.ID),
	)
	if err != nil {
		return nil, xerror.NewFatalError("error to libp2p.New").Wrap(err)
	}
	cc.Host = host

	pgpDecrypter, err := crypto.NewPGPDecrypter(cfg.PGPPrivateKey)
	if err != nil {
		return nil, xerror.NewFatalError("error to crypto.NewPGPDecrypter").Wrap(err)
	}

	go func() {
		err := cc.DeviceManager.Load()
		if err != nil {
			cc.ErrorChan <- xerror.NewFatalError("can not load device from setting").Wrap(err)
		}

		streamHandler := stream.NewStreamHandler(
			cc.Config,
			cc.ClipboardManager,
			cc.DeviceManager,
			cc.LogChan,
			cc.ErrorChan,
			pgpDecrypter,
		)
		cc.streamHandler = streamHandler

		// This function is called when a peer initiates a connection and starts a stream with this peer.
		cc.Host.SetStreamHandler(stream.PROTOCAL_ID, streamHandler.HandleStream)
		cc.LogChan <- fmt.Sprintf("[*] your multiaddress is: /ip4/%s/tcp/%v/p2p/%s", cc.Config.ListenHost, cc.Config.ListenPort, host.ID().Pretty())

		peerInfoChan, err := discovery.InitMultiMDNS(cc.Host, cc.Config.GroupName, cc.LogChan)
		if err != nil {
			cc.ErrorChan <- xerror.NewFatalError("error to discovery.InitMultiMDNS").Wrap(err)
		}

		for peerInfo := range peerInfoChan { // when discover a peer
			dv := cc.DeviceManager.GetDevice(peerInfo.ID.Pretty())
			if dv != nil && dv.Status == device.StatusBlocked {
				cc.ErrorChan <- xerror.NewRuntimeErrorf("device %s is blocked", peerInfo.ID.Pretty())
				continue
			}

			cc.LogChan <- fmt.Sprintf("connecting to peer host: %s", peerInfo)

			if err := cc.Host.Connect(ctx, peerInfo); err != nil {
				cc.ErrorChan <- xerror.NewRuntimeError("connect error").Wrap(err)
				continue
			}

			// open a stream, this stream will be handled by handleStream other end
			stream, err := cc.Host.NewStream(ctx, peerInfo.ID, stream.PROTOCAL_ID)
			if err != nil {
				cc.ErrorChan <- xerror.NewRuntimeError("new stream error").Wrap(err)
				continue
			}

			if dv == nil {
				dv = device.NewDevice(peerInfo, stream)
			} else {
				dv.AddressInfo = peerInfo
				dv.Stream = stream
				dv.Reader = bufio.NewReader(stream)
				dv.Writer = bufio.NewWriter(stream)
			}

			cc.DeviceManager.UpdateDevice(dv)
			go streamHandler.CreateReadData(dv.Reader, dv)

			cc.LogChan <- fmt.Sprintf("connected to peer host: %s", peerInfo)
		}
	}()

	return cc, nil
}

func (cc *CrossClipboard) Stop() error {
	if cc.streamHandler != nil {
		for id, dv := range cc.DeviceManager.Devices {
			// graceful close connection stream
			if dv.Status == device.StatusConnected {
				log.Printf("ending stream for peer %s \n", id)

				cc.streamHandler.SendSignal(dv, stream.SignalDisconnect)

				dv.Stream.Close()
			}
		}
	}

	cc = nil
	return nil
}

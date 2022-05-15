package p2p

import (
	"bufio"
	"context"
	"crypto/rand"
	"fmt"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/multiformats/go-multiaddr"
	"github.com/ntsd/cross-clipboard/pkg/clipboard"
	"github.com/ntsd/cross-clipboard/pkg/discovery"
	"github.com/ntsd/cross-clipboard/pkg/stream"
	"github.com/ntsd/cross-clipboard/pkg/utils"
)

func StartP2P(cfg utils.Config) {
	ctx := context.Background()

	// Creates a new ECDSA key pair for this host.
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.ECDSA, 2048, rand.Reader)
	if err != nil {
		panic(err)
	}

	// 0.0.0.0 will listen on any interface device.
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", cfg.ListenHost, cfg.ListenPort))
	fmt.Printf("[*] Listening on: %s with port: %d\n", cfg.ListenHost, cfg.ListenPort)

	// libp2p.New constructs a new libp2p Host.
	host, err := libp2p.New(
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(prvKey),
	)
	if err != nil {
		panic(err)
	}

	// initial clipboard
	cb := clipboard.NewClipboard()
	streamHandler := stream.NewStreamHandler(cb)

	// Set a function as stream handler.
	// This function is called when a peer initiates a connection and starts a stream with this peer.
	host.SetStreamHandler(protocol.ID(cfg.ProtocolID), streamHandler.HandleStream)

	fmt.Printf("\n[*] Your Multiaddress Is: /ip4/%s/tcp/%v/p2p/%s\n", cfg.ListenHost, cfg.ListenPort, host.ID().Pretty())

	peerChan := discovery.InitMultiMDNS(host, cfg.RendezvousString)

	for peer := range peerChan { // when discover a peer
		fmt.Println("Found peer:", peer, ", connecting")

		if err := host.Connect(ctx, peer); err != nil {
			fmt.Println("Connection failed:", err)
			continue
		}

		// open a stream, this stream will be handled by handleStream other end
		stream, err := host.NewStream(ctx, peer.ID, protocol.ID(cfg.ProtocolID))
		if err != nil {
			fmt.Println("Stream open failed", err)
			continue
		}
		// Add reader
		go streamHandler.CreateReadData(bufio.NewReader(stream), peer.ID.Pretty())
		// Add writer
		streamHandler.AddWriter(bufio.NewWriter(stream), peer.ID.Pretty())
		fmt.Println("Connected to:", peer)
	}

	select {} //wait here
}

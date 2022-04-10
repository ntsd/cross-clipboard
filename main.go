package main

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"log"
	mrand "math/rand"
	"os"

	"github.com/ntsd/cross-clipboard/p2p"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sourcePort := flag.Int("sp", 0, "Source port number")
	dest := flag.String("d", "", "Destination multiaddr string")
	help := flag.Bool("help", false, "Display help")
	debug := flag.Bool("debug", false, "Debug generates the same node ID on every execution")

	flag.Parse()

	if *help {
		fmt.Printf("This program demonstrates a simple p2p chat application using libp2p\n\n")
		fmt.Println("Usage: Run './chat -sp <SOURCE_PORT>' where <SOURCE_PORT> can be any port number.")
		fmt.Println("Now run './chat -d <MULTIADDR>' where <MULTIADDR> is multiaddress of previous listener host.")

		os.Exit(0)
	}

	// If debug is enabled, use a constant random source to generate the peer ID. Only useful for debugging
	var r io.Reader
	if *debug {
		r = mrand.New(mrand.NewSource(int64(*sourcePort)))
	} else {
		r = rand.Reader
	}

	h, err := p2p.MakeHost(*sourcePort, r)
	if err != nil {
		log.Println(err)
		return
	}

	if *dest == "" {
		p2p.StartPeer(ctx, h, p2p.HandleStream)
	} else {
		rw, err := p2p.StartPeerAndConnect(ctx, h, *dest)
		if err != nil {
			log.Println(err)
			return
		}

		// Create a thread to read and write data.
		go p2p.WriteData(rw)
		go p2p.ReadData(rw)
	}

	// Wait forever
	select {}

	// help := flag.Bool("help", false, "Display Help")
	// cfg := utils.ParseFlags()

	// if *help {
	// 	fmt.Printf("Simple example for peer discovery using mDNS. mDNS is great when you have multiple peers in local LAN.")
	// 	fmt.Printf("Usage: \n   Run './chat-with-mdns'\nor Run './chat-with-mdns -host [host] -port [port] -rendezvous [string] -pid [proto ID]'\n")

	// 	os.Exit(0)
	// }

	// fmt.Printf("[*] Listening on: %s with port: %d\n", cfg.ListenHost, cfg.ListenPort)

	// ctx := context.Background()
	// r := rand.Reader

	// // Creates a new RSA key pair for this host.
	// prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	// if err != nil {
	// 	panic(err)
	// }

	// // 0.0.0.0 will listen on any interface device.
	// sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", cfg.ListenHost, cfg.ListenPort))

	// // libp2p.New constructs a new libp2p Host.
	// // Other options can be added here.
	// host, err := libp2p.New(
	// 	libp2p.ListenAddrs(sourceMultiAddr),
	// 	libp2p.Identity(prvKey),
	// )
	// if err != nil {
	// 	panic(err)
	// }

	// // Set a function as stream handler.
	// // This function is called when a peer initiates a connection and starts a stream with this peer.
	// host.SetStreamHandler(protocol.ID(cfg.ProtocolID), p2p.HandleStream)

	// fmt.Printf("\n[*] Your Multiaddress Is: /ip4/%s/tcp/%v/p2p/%s\n", cfg.ListenHost, cfg.ListenPort, host.ID().Pretty())

	// discoveryNotifees := p2p.InitMultiMDNS(host, cfg.RendezvousString)

	// for id, discoveryNotifee := range discoveryNotifees {
	// 	if id == host.ID().Pretty() {
	// 		continue
	// 	}

	// 	fmt.Println("id", id)

	// 	peer := <-discoveryNotifee.PeerChan // will block untill we discover a peer
	// 	fmt.Println("Found peer:", peer, ", connecting")

	// 	if err := host.Connect(ctx, peer); err != nil {
	// 		fmt.Println("Connection failed:", err)
	// 	}

	// 	// open a stream, this stream will be handled by handleStream other end
	// 	stream, err := host.NewStream(ctx, peer.ID, protocol.ID(cfg.ProtocolID))

	// 	if err != nil {
	// 		fmt.Println("Stream open failed", err)
	// 	} else {
	// 		rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	// 		go p2p.WriteData(rw)
	// 		go p2p.ReadData(rw)
	// 		fmt.Println("Connected to:", peer)
	// 	}
	// }

	// select {} //wait here
}

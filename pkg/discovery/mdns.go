package discovery

import (
	"fmt"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

type DiscoveryNotifee struct {
	PeerHost host.Host
	PeerChan chan peer.AddrInfo
}

//interface to be called when new  peer is found
func (n *DiscoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	fmt.Println("Discovered", pi)
	if n.PeerHost.ID().Pretty() != pi.ID.Pretty() {
		n.PeerChan <- pi
	}
}

//Initialize the MDNS service
func InitMultiMDNS(peerhost host.Host, rendezvous string) chan peer.AddrInfo {
	// register with service so that we get notified about peer discovery
	n := &DiscoveryNotifee{
		PeerHost: peerhost,
	}
	n.PeerChan = make(chan peer.AddrInfo)

	// An hour might be a long long period in practical applications. But this is fine for us
	ser := mdns.NewMdnsService(peerhost, rendezvous, n)
	if err := ser.Start(); err != nil {
		panic(err)
	}

	return n.PeerChan
}

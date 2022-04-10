package p2p

import (
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

type DiscoveryNotifee struct {
	PeerChan chan peer.AddrInfo
}

//interface to be called when new  peer is found
func (n *DiscoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	n.PeerChan <- pi
}

//Initialize the MDNS service
func InitMultiMDNS(peerhost host.Host, rendezvous string) map[string]*DiscoveryNotifee {

	maxPeerNum := 100
	discoveryNotifees := map[string]*DiscoveryNotifee{}

	for i := 0; i < maxPeerNum; i++ {
		// register with service so that we get notified about peer discovery
		n := &DiscoveryNotifee{}
		n.PeerChan = make(chan peer.AddrInfo)

		// An hour might be a long long period in practical applications. But this is fine for us
		ser := mdns.NewMdnsService(peerhost, rendezvous, n)
		if err := ser.Start(); err != nil {
			panic(err)
		}

		peer := <-n.PeerChan

		discoveryNotifees[peer.ID.Pretty()] = n
	}

	return discoveryNotifees
}

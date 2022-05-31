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
	LogChan  chan string
}

//interface to be called when new  peer is found
func (n *DiscoveryNotifee) HandlePeerFound(peerInfo peer.AddrInfo) {
	n.LogChan <- fmt.Sprintf("discovered: %s", peerInfo)
	if n.PeerHost.ID().Pretty() != peerInfo.ID.Pretty() {
		n.PeerChan <- peerInfo
	}
}

//Initialize the MDNS service
func InitMultiMDNS(peerhost host.Host, rendezvous string, logchan chan string) (chan peer.AddrInfo, error) {
	// register with service so that we get notified about peer discovery
	n := &DiscoveryNotifee{
		PeerHost: peerhost,
		PeerChan: make(chan peer.AddrInfo),
		LogChan:  logchan,
	}

	// An hour might be a long long period in practical applications. But this is fine for us
	ser := mdns.NewMdnsService(peerhost, rendezvous, n)
	if err := ser.Start(); err != nil {
		return nil, err
	}

	return n.PeerChan, nil
}

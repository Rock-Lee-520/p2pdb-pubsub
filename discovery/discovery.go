package discovery

import (
	"context"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

// DiscoveryInterval is how often we re-publish our mDNS records.
const DiscoveryInterval = time.Hour

// DiscoveryServiceTag is used in our mDNS advertisements to discover other chat peers.
const DiscoveryServiceTag = "p2pdb-example"

type Discovery interface {
	// create a new libp2p Host that listens on a random TCP port
	//	Connect() (Host string)
	Create(host string) (host.Host, error)
	SetupDiscovery(host host.Host) error
}

type DiscoveryFactory struct {
}

const LISTEN_ADDRESS_STRINGS = "/ip4/0.0.0.0/tcp/0"

// discoveryNotifee gets notified when we find a new peer via mDNS discovery
type discoveryNotifee struct {
	h host.Host
}

func NewDiscoveryFactory() *DiscoveryFactory {
	return &DiscoveryFactory{}
}

// func (d *DiscoveryFactory) Connect() (Host string)  (Discovery, error){
// 	return d,nil
// }

func (d *DiscoveryFactory) Create(host string) (host.Host, error) {
	if host == "" {
		host = LISTEN_ADDRESS_STRINGS
	}
	h, err := libp2p.New(libp2p.ListenAddrStrings(host))

	return h, err
}

// setupDiscovery creates an mDNS discovery service and attaches it to the libp2p Host.
// This lets us automatically discover peers on the same LAN and connect to them.
func (d *DiscoveryFactory) SetupDiscovery(h host.Host) error {
	// setup mDNS discovery to find local peers
	s := mdns.NewMdnsService(h, DiscoveryServiceTag, &discoveryNotifee{h: h})
	return s.Start()
}

// HandlePeerFound connects to peers discovered via mDNS. Once they're connected,
// the PubSub system will automatically start interacting with them if they also
// support PubSub.
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	fmt.Printf("discovered new peer %s\n", pi.ID.Pretty())
	err := n.h.Connect(context.Background(), pi)
	if err != nil {
		fmt.Printf("error connecting to peer %s: %s\n", pi.ID.Pretty(), err)
	}
}

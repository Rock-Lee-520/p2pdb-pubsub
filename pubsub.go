package pubsub

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-core/peer"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/mr-tron/base58/base58"
	"github.com/multiformats/go-multihash"
	"github.com/sirupsen/logrus"
	"time"

	"github.com/Rock-liyi/p2pdb/infrastructure/util/log"
	"github.com/libp2p/go-libp2p-core/host"
	libpubsub "github.com/libp2p/go-libp2p-pubsub"
)

type SubHandler interface {
	SubHandler(DataMessage)
}

type PubSubInterface interface {
	InitPubSub()
	GetType() string
	SetType(Type string)
	Pub(message DataMessage)
	Sub() (*libpubsub.Subscription, *libpubsub.Topic, error)
}

type PubSub struct {
	Host       host.Host
	Topic      *libpubsub.Topic
	Type       string
	SubHandler SubHandler
	ctx        context.Context
}

type DataMessage struct {
	Type string
	Data interface{}
}

const GlobalTopic string = ("p2pdb")

func (p2pdbPubSub *PubSub) GetType() string {
	return p2pdbPubSub.Type
}

func (p2pdbPubSub *PubSub) SetType(Type string) {
	p2pdbPubSub.Type = Type
}

func (p2pdbPubSub *PubSub) SetHost(Host host.Host) {
	p2pdbPubSub.Host = Host
}

type discoveryNotifee struct {
	h host.Host
}

func (d discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	fmt.Printf("discovered new peer %s\n", pi.ID.Pretty())
	err := d.h.Connect(context.Background(), pi)
	if err != nil {
		fmt.Printf("error connecting to peer %s: %s\n", pi.ID.Pretty(), err)
	}
}

func (p2pdbPubSub *PubSub) InitPubSub(ctx context.Context, Type string, Host host.Host, Routingdiscovery *discovery.RoutingDiscovery, dht *dht.IpfsDHT) PubSub {

	p2pdbPubSub.ctx = ctx
	p2pdbPubSub.Host = Host
	p2pdbPubSub.Type = Type

	ps, err := libpubsub.NewGossipSub(p2pdbPubSub.ctx, p2pdbPubSub.Host, libpubsub.WithDiscovery(Routingdiscovery))
	if err != nil {
		panic(err)
	}
	if p2pdbPubSub.Type == "" {
		p2pdbPubSub.Type = GlobalTopic
	}

	s := mdns.NewMdnsService(p2pdbPubSub.Host, GlobalTopic, &discoveryNotifee{h: p2pdbPubSub.Host})
	err = s.Start()
	if err != nil {
		//panic(err)
		//panic(err)
		log.Error(err.Error())
	}

	//p2pdbPubSub.ConnectPeers(ctx, dht)

	topic, err := ps.Join(p2pdbPubSub.Type)
	if err != nil {
		panic(err)
	}
	p2pdbPubSub.Topic = topic

	// and subscribe to it
	sub, err := topic.Subscribe()
	time.Sleep(3 * time.Second)

	go p2pdbPubSub.startNewSubscribeService(sub)
	return *p2pdbPubSub
}

func (p2pdbPubSub *PubSub) Pub(message DataMessage) {
	log.Debug("call method PubSub.Pub")
	log.Debug(message.Type)
	log.Debug(p2pdbPubSub.GetType())
	data, err := json.Marshal(message)
	if err != nil {
		log.Panic(err)
	}
	err = p2pdbPubSub.Topic.Publish(p2pdbPubSub.ctx, data)
	if err != nil {
		log.Panic(err)
	}
}

func (p2pdbPubSub *PubSub) startNewSubscribeService(sub *libpubsub.Subscription) {
	for {
		msg, err := sub.Next(p2pdbPubSub.ctx)
		if err != nil {
			fmt.Printf("error:%s\n", err.Error())
			//panic(err)
			return
		}

		// only forward messages delivered by others
		cm := new(DataMessage)
		err = json.Unmarshal(msg.Data, cm)
		if err != nil {
			continue
		}
		log.Info("收到消息啦=======")

		log.Debug(cm)
		//recieve messages to other channels
		// var newMessage event.Message
		// newMessage.Data = function.JsonEncode(cm.Data)
		// newMessage.Type = cm.Type

		// event.PublishSyncEvent(cm.Type, newMessage)
		//p2pdbPubSub.SubHandler(cm)
	}
}

func (p2pdbPubSub *PubSub) ConnectPeers(ctx context.Context, dht *dht.IpfsDHT) {
	// Generate the Service CID
	cidvalue := generateCID(p2pdbPubSub.Host.ID().String())

	err := dht.Provide(ctx, cidvalue, true)
	if err != nil {
		log.Debug("dht Provide fail", err.Error())
	}

	time.Sleep(time.Second * 5)

	peerchan := dht.FindProvidersAsync(ctx, cidvalue, 0)
	go handlePeerDiscovery(p2pdbPubSub.Host, peerchan)

}

func handlePeerDiscovery(h host.Host, peerchan <-chan peer.AddrInfo) {
	// Iterate over the peer channel
	for peer := range peerchan {
		// Ignore if the discovered peer is the host itself
		if peer.ID == h.ID() {
			continue
		}

		// Connect to the peer
		err := h.Connect(context.Background(), peer)

		if err != nil {
			log.Debug("peer connection failed: ", err.Error())
		}

		log.Debug("peer connection success: ", peer.ID)
	}
}

// A function that generates a CID object for a given string and returns it.
// Uses SHA256 to hash the string and generate a multihash from it.
// The mulithash is then base58 encoded and then used to create the CID
func generateCID(namestring string) cid.Cid {
	// Hash the service content ID with SHA256
	hash := sha256.Sum256([]byte(namestring))
	// Append the hash with the hashing codec ID for SHA2-256 (0x12),
	// the digest size (0x20) and the hash of the service content ID
	finalhash := append([]byte{0x12, 0x20}, hash[:]...)
	// Encode the fullhash to Base58
	b58string := base58.Encode(finalhash)

	// Generate a Multihash from the base58 string
	mulhash, err := multihash.FromB58String(string(b58string))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatalln("Failed to Generate Service CID!")
	}

	// Generate a CID from the Multihash
	cidvalue := cid.NewCidV1(12, mulhash)
	// Return the CID
	return cidvalue
}

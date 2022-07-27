package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	discovery "github.com/Rock-liyi/p2pdb-discovery"
	"github.com/Rock-liyi/p2pdb/application/event"
	"github.com/libp2p/go-libp2p-core/host"
	libpubsub "github.com/libp2p/go-libp2p-pubsub"
)

// DataMessage gets converted to/from JSON and sent in the body of pubsub messages.
type DataMessage struct {
	data      interface{}
	topic     string
	requestId string
}

type SubHandler interface {
	SubHandler(DataMessage)
}

type PubSubInfterface interface {
	InitPub()
	GetType() string
	SetType(Type string)
	Pub(message DataMessage)
	initDiscovery()
}

type PubSub struct {
	h          host.Host
	topic      *libpubsub.Topic
	Type       string
	SubHandler SubHandler
	ctx        context.Context
	discovery  discovery.Discovery
}

const Topic string = ("p2pdb")
const Address string = "/ip4/0.0.0.0/tcp/0"

func (p2pdbPubSub *PubSub) GetType() string {
	return p2pdbPubSub.Type
}

func (p2pdbPubSub *PubSub) SetType(Type string) {
	p2pdbPubSub.Type = Type
}

func (p2pdbPubSub *PubSub) InitPub() {
	p2pdbPubSub.ctx = context.Background()
	p2pdbPubSub.initDiscovery()

	// create a new PubSub service using the GossipSub router
	ps, err := libpubsub.NewGossipSub(p2pdbPubSub.ctx, p2pdbPubSub.h)
	if err != nil {
		panic(err)
	}
	if p2pdbPubSub.Type == "" {
		p2pdbPubSub.Type = Topic
	}

	topic, err := ps.Join(p2pdbPubSub.Type)
	if err != nil {
		panic(err)
	}
	p2pdbPubSub.topic = topic
	time.Sleep(1 * time.Second)
}

func (p2pdbPubSub *PubSub) Pub(message DataMessage) {
	data, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	p2pdbPubSub.topic.Publish(p2pdbPubSub.ctx, data)
}

func (p2pdbPubSub *PubSub) Sub() {

	p2pdbPubSub.ctx = context.Background()

	p2pdbPubSub.initDiscovery()

	// create a new PubSub service using the GossipSub router
	ps, err := libpubsub.NewGossipSub(p2pdbPubSub.ctx, p2pdbPubSub.h)
	if err != nil {
		panic(err)
	}

	topic, err := ps.Join(Topic)
	if err != nil {
		panic(err)
	}

	// and subscribe to it
	sub, err := topic.Subscribe()
	if err != nil {
		panic(err)
	}

	for {
		msg, err := sub.Next(p2pdbPubSub.ctx)
		if err != nil {
			fmt.Printf("error:%s\n", err.Error())
			panic(err)
			return
		}
		// only forward messages delivered by others
		cm := new(DataMessage)
		err = json.Unmarshal(msg.Data, cm)
		if err != nil {
			continue
		}

		//recieve messages to other channels
		var newMessage event.Message
		newMessage.Data = cm.data
		newMessage.Type = cm.topic
		event.PublishSyncEvent(cm.topic, newMessage)
		//p2pdbPubSub.SubHandler(cm)
	}
}

// initDiscovery
func (p2pdbPubSub *PubSub) initDiscovery() {
	p2pdbPubSub.discovery = discovery.NewDiscoveryFactory()
	// create a new libp2p Host that listens on a random TCP port
	h, err := p2pdbPubSub.discovery.Create(Address)
	if err != nil {
		panic(err)
	}
	p2pdbPubSub.h = h

	// setup local mDNS discovery
	if err := p2pdbPubSub.discovery.SetupDiscovery(p2pdbPubSub.h); err != nil {
		panic(err)
	}
}

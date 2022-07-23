package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	discovery "github.com/Rock-liyi/p2pdb-discovery"
	"github.com/libp2p/go-libp2p-core/host"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"time"
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

type PubSub struct {
	h          host.Host
	topic      *pubsub.Topic
	SubHandler SubHandler
	ctx        context.Context
	discovery  discovery.Discovery
}

const Topic string = ("p2pDB")
const Address string = "/ip4/0.0.0.0/tcp/0"

func (p2pdbPubSub *PubSub) InitPublish() {
	p2pdbPubSub.ctx = context.Background()
	p2pdbPubSub.initDiscovery()

	// create a new PubSub service using the GossipSub router
	ps, err := pubsub.NewGossipSub(p2pdbPubSub.ctx, p2pdbPubSub.h)
	if err != nil {
		panic(err)
	}

	topic, err := ps.Join(Topic)
	if err != nil {
		panic(err)
	}
	p2pdbPubSub.topic = topic
	time.Sleep(1 * time.Second)
}

func (p2pdbPubSub *PubSub) Publish(message DataMessage) {

	if p2pdbPubSub.topic == nil {
		p2pdbPubSub.InitPublish()
	}

	data, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	p2pdbPubSub.topic.Publish(p2pdbPubSub.ctx, data)
}

func (p2pdbPubSub *PubSub) Subscribe() {

	p2pdbPubSub.ctx = context.Background()

	p2pdbPubSub.initDiscovery()

	// create a new PubSub service using the GossipSub router
	ps, err := pubsub.NewGossipSub(p2pdbPubSub.ctx, p2pdbPubSub.h)
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
		p2pdbPubSub.SubHandler(cm)
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

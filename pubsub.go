package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	discovery "github.com/Rock-liyi/p2pdb-discovery"
	"github.com/Rock-liyi/p2pdb/application/event"
	"github.com/Rock-liyi/p2pdb/infrastructure/util/function"
	"github.com/Rock-liyi/p2pdb/infrastructure/util/log"
	"github.com/libp2p/go-libp2p-core/host"
	libpubsub "github.com/libp2p/go-libp2p-pubsub"
)

type SubHandler interface {
	SubHandler(DataMessage)
}

type PubSubInterface interface {
	InitPub()
	GetType() string
	SetType(Type string)
	Pub(message DataMessage)
	Sub() (*libpubsub.Subscription, *libpubsub.Topic, error)
	initDiscovery()
}

type PubSub struct {
	Host       host.Host
	Topic      *libpubsub.Topic
	Type       string
	SubHandler SubHandler
	ctx        context.Context
	discovery  discovery.Discovery
	Message    DataMessage
}

const GlobalTopic string = ("p2pdb")
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
	ps, err := libpubsub.NewGossipSub(p2pdbPubSub.ctx, p2pdbPubSub.Host)
	if err != nil {
		panic(err)
	}
	if p2pdbPubSub.Type == "" {
		p2pdbPubSub.Type = GlobalTopic
	}

	topic, err := ps.Join(p2pdbPubSub.Type)
	if err != nil {
		panic(err)
	}
	p2pdbPubSub.Topic = topic
	time.Sleep(5 * time.Second)
}

func (p2pdbPubSub *PubSub) Pub(message DataMessage) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Panic(err)
	}
	p2pdbPubSub.Topic.Publish(p2pdbPubSub.ctx, data)
}

func (p2pdbPubSub *PubSub) Sub() (*libpubsub.Subscription, *libpubsub.Topic, error) {

	p2pdbPubSub.ctx = context.Background()

	p2pdbPubSub.initDiscovery()

	// create a new PubSub service using the GossipSub router
	ps, err := libpubsub.NewGossipSub(p2pdbPubSub.ctx, p2pdbPubSub.Host)
	if err != nil {
		panic(err)
	}

	log.Debug("topic is " + GlobalTopic)
	topic, err := ps.Join(GlobalTopic)
	if err != nil {
		return nil, nil, err
	}

	// and subscribe to it
	sub, err := topic.Subscribe()
	if err != nil {
		return nil, nil, err
	}
	return sub, topic, nil
}

func (p2pdbPubSub *PubSub) StartNewSubscribeService(sub *libpubsub.Subscription) {
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
		newMessage.Data = function.JsonEncode(cm.Data)
		newMessage.Type = cm.Type

		event.PublishSyncEvent(cm.Type, newMessage)
		//p2pdbPubSub.SubHandler(cm)
	}
}

// DataMessage gets converted to/from JSON and sent in the body of pubsub messages.
// type DataMessage struct {
// 	data      interface{}
// 	topic     string
// 	requestId string
// }

type DataMessage struct {
	Type string
	Data interface{}
}

// initDiscovery
func (p2pdbPubSub *PubSub) initDiscovery() {
	if p2pdbPubSub.discovery != nil {
		return
	}
	p2pdbPubSub.discovery = discovery.NewDiscoveryFactory()
	// create a new libp2p Host that listens on a random TCP port
	h, err := p2pdbPubSub.discovery.Create(Address)
	if err != nil {
		panic(err)
	}
	p2pdbPubSub.Host = h

	// setup local mDNS discovery
	if err := p2pdbPubSub.discovery.SetupDiscovery(p2pdbPubSub.Host); err != nil {
		panic(err)
	}
}

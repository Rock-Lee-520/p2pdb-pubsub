package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/libp2p/go-libp2p-core/peer"
	discovery "github.com/libp2p/go-libp2p-discovery"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
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

func (p2pdbPubSub *PubSub) InitPubSub(ctx context.Context, Type string, Host host.Host, Routingdiscovery *discovery.RoutingDiscovery) PubSub {

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
		panic(err)
	}

	topic, err := ps.Join(p2pdbPubSub.Type)
	if err != nil {
		panic(err)
	}
	p2pdbPubSub.Topic = topic

	// and subscribe to it
	sub, err := topic.Subscribe()
	time.Sleep(3 * time.Second)

	go p2pdbPubSub.StartNewSubscribeService(sub)
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

func (p2pdbPubSub *PubSub) StartNewSubscribeService(sub *libpubsub.Subscription) {
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

// DataMessage gets converted to/from JSON and sent in the body of pubsub messages.
// type DataMessage struct {
// 	data      interface{}
// 	topic     string
// 	requestId string
// }

// initDiscovery
//func (p2pdbPubSub *PubSub) initDiscovery() {
//	if p2pdbPubSub.discovery != nil {
//		return
//	}
//	p2pdbPubSub.discovery = discovery.NewDiscoveryFactory()
//
//	// create a new libp2p Host that listens on a random TCP port
//	h, err := p2pdbPubSub.discovery.Create(Address)
//	if err != nil {
//		panic(err)
//	}
//	p2pdbPubSub.Host = h
//
//	// setup local mDNS discovery
//	if err := p2pdbPubSub.discovery.SetupDiscovery(p2pdbPubSub.Host); err != nil {
//		panic(err)
//	}
//}

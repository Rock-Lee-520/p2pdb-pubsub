package main

import (
	"context"
	"encoding/json"
	"time"

	discovery "github.com/Rock-liyi/p2pdb-discovery"

	debug "github.com/favframework/debug"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

func main() {
	ctx := context.Background()
	Discovery := discovery.NewDiscoveryFactory()
	// create a new libp2p Host that listens on a random TCP port
	h, err := Discovery.Create("/ip4/0.0.0.0/tcp/0")
	debug.Dump(h.ID().String())
	// create a new PubSub service using the GossipSub router
	ps, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		panic(err)
	}

	// setup local mDNS discovery
	if err := Discovery.SetupDiscovery(h); err != nil {
		panic(err)
	}

	topic, err := ps.Join((topicName("p2pdb")))
	if err != nil {
		panic(err)
	}

	// and subscribe to it
	sub, err := topic.Subscribe()
	if err != nil {
		panic(err)
	}

	for {
		time.Sleep(1 * time.Second)
		msg, err := sub.Next(ctx)
		if err != nil {
			panic(err)
			return
		}
		// only forward messages delivered by others

		cm := new(DataMessage)
		err = json.Unmarshal(msg.Data, cm)
		debug.Dump("debug:")
		debug.Dump(cm.Message)
		if err != nil {
			continue
		}
		debug.Dump(cm)
		// send valid messages onto the Messages channel
		//Messages <- cm
	}
}

// ChatMessage gets converted to/from JSON and sent in the body of pubsub messages.
type DataMessage struct {
	Message    string
	SenderID   string
	SenderNick string
}

func topicName(Name string) string {
	return "chat-room:" + Name
}

// // HandlePeerFound connects to peers discovered via mDNS. Once they're connected,
// // the PubSub system will automatically start interacting with them if they also
// // support PubSub.
// func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
// 	fmt.Printf("discovered new peer %s\n", pi.ID.Pretty())
// 	err := n.h.Connect(context.Background(), pi)
// 	if err != nil {
// 		fmt.Printf("error connecting to peer %s: %s\n", pi.ID.Pretty(), err)
// 	}
// }

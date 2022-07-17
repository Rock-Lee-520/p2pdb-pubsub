package main

import (
	"context"
	"encoding/json"
	"time"

	discovery "github.com/Rock-liyi/p2pdb-discovery"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// DiscoveryInterval is how often we re-publish our mDNS records.
// const DiscoveryInterval = time.Hour

// // DiscoveryServiceTag is used in our mDNS advertisements to discover other chat peers.
// const DiscoveryServiceTag = "p2pdb-example"

func main() {
	ctx := context.Background()

	Discovery := discovery.NewDiscoveryFactory()
	//create a new libp2p Host that listens on a random TCP port
	// h, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
	// if err != nil {
	// 	panic(err)
	// }

	h, err := Discovery.Create("/ip4/0.0.0.0/tcp/0")

	if err != nil {
		panic(err)
	}
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

	for {
		time.Sleep(1 * time.Second)
		input := DataMessage{
			Message:    "123",
			SenderID:   "123",
			SenderNick: "1233",
		}

		msgBytes, err := json.Marshal(input)
		if err != nil {
			panic(err)
		}
		topic.Publish(ctx, msgBytes)
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

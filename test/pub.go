package main

import (
	"time"

	pubsub "github.com/Rock-liyi/p2pdb-pubsub"
)

func main() {
	var Sub = &pubsub.PubSub{}
	// Sub.SetType("p2pdb")
	// var subscription, err = Sub.Sub()
	// if err != nil {
	// 	panic(err)
	// }
	Sub.InitPub()
	Sub.Pub(pubsub.DataMessage{Type: "test", Data: "TestPubsub"})

	time.Sleep(5 * time.Second)
	//Sub.StartNewSubscribeService(subscription)
}

package pubsub

import (
	"testing"
	"time"
)

func TestPubsub(t *testing.T) {
	var Sub = &PubSub{}
	// Sub.SetType("p2pdb")
	// var subscription, err = Sub.Sub()
	// if err != nil {
	// 	panic(err)
	// }
	Sub.InitPub()
	Sub.Pub(DataMessage{Type: "test", Data: "TestPubsub"})

	time.Sleep(5 * time.Second)
	//Sub.StartNewSubscribeService(subscription)
}

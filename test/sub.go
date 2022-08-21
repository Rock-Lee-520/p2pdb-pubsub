package main

import (
	pubsub "github.com/Rock-liyi/p2pdb-pubsub"
)

func main() {
	var Sub = &pubsub.PubSub{}
	var topicType = "p2pdb"
	Sub.SetType(topicType)
	//var subscription, _, err = Sub.Sub()
	//if err != nil {
	//	panic(err)
	//}

	//Sub.StartNewSubscribeService(subscription)
}

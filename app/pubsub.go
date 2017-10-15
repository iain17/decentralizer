package app

import (
	"context"
	"strings"
	"github.com/ipfs/go-ipfs/core"
	"fmt"
	peer "gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
)

func Receive(node *core.IpfsNode, topic string, didChange func(peer peer.ID, message string)) {
	sub, err := node.Floodsub.Subscribe(topic)
	if err != nil {
		panic(err)
	}

	var lastMessage string
	for {
		logger.Debug("looking for messages")
		msg, err := sub.Next(context.Background())
		if err != nil {
			fmt.Println(err)
			continue
		}
		message := strings.TrimSpace(string(msg.GetData()))
		peer := msg.GetFrom()

		if lastMessage == message {
			continue
		}
		didChange(peer, message)

		//Send to the nodes I'm connected to.
		Publish(node, topic, message)
		lastMessage = message
	}
}

func Publish(node *core.IpfsNode, topic string, message string) {
	//st := fmt.Sprintf("%s::#%s", sp.String(), h.String())
	fmt.Println("Publishing", message)
	node.Floodsub.Publish(topic, []byte(message))
}
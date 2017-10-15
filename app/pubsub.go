package app

import (
	"context"
	"strings"
	"fmt"
	"gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
)

func (d *Decentralizer) Receive(topic string, didChange func(peer peer.ID, message string)) {
	sub, err := d.i.Floodsub.Subscribe(topic)
	if err != nil {
		panic(err)
	}

	var lastMessage string
	for {
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
		if peer.String() != d.i.Identity.String() {
			didChange(peer, message)
		}

		//Send to the nodes I'm connected to.
		d.Publish(topic, message)
		lastMessage = message
	}
}

func (d *Decentralizer) Publish(topic string, message string) {
	//st := fmt.Sprintf("%s::#%s", sp.String(), h.String())
	fmt.Println("Publishing", message)
	test := d.i.Peerstore.Peers()
	println(len(test))
	println(test)
	d.i.Floodsub.Publish(topic, []byte(message))
}
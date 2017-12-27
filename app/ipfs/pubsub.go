package ipfs

import (
	"context"
	"github.com/iain17/logger"
	"gx/ipfs/QmYHpXQEWuhwgRFBnrf4Ua6AZhcqXCYa7Biv65SLGgTgq5/go-ipfs/core"
	"gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
	"gx/ipfs/QmUUSLfvihARhCxxgnjW4hmycJpPvzNu12Aaz6JWVdfnLg/go-libp2p-floodsub"
	"io"
)

func Subscribe(node *core.IpfsNode, topic string, didChange func(peer peer.ID, data []byte)) (*floodsub.Subscription, error) {
	sub, err := node.Floodsub.Subscribe(topic)
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			msg, err := sub.Next(context.Background())
			if err == io.EOF || err == context.Canceled {
				return
			} else if err != nil {
				logger.Error(err)
				return
			}
			peer := msg.GetFrom()
			if peer.String() != node.Identity.String() {
				didChange(peer, msg.GetData())
			}
		}
	}()
	return sub, err
}

func Publish(node *core.IpfsNode, topic string, data []byte) error {
	return node.Floodsub.Publish(topic, data)
}

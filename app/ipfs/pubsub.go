package ipfs

import (
	"context"
	"github.com/iain17/logger"
	"gx/ipfs/QmNUKMfTHQQpEwE8bUdv5qmKC3ymdW7zw82LFS8D6MQXmu/go-ipfs/core"
	"gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
	"gx/ipfs/Qmdnza7rLi7CMNNwNhNkcs9piX5sf6rxE8FrCsPzYtUEUi/floodsub"
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

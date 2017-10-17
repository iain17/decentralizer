package ipfs

import (
	"context"
	"github.com/ipfs/go-ipfs/core"
	"gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
	"github.com/iain17/logger"
	"gx/ipfs/QmTm7GoSkSSQPP32bZhvu17oY1AfvPKND6ELUdYAcKuR1j/floodsub"
)

func Subscribe(node *core.IpfsNode, topic string, didChange func(peer peer.ID, data []byte)) (*floodsub.Subscription, error) {
	sub, err := node.Floodsub.Subscribe(topic)
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			msg, err := sub.Next(context.Background())
			if err != nil {
				logger.Debug(err)
				break
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
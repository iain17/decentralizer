package ipfs

import (
	"context"
	"github.com/iain17/logger"
	"github.com/pkg/errors"
	"gx/ipfs/QmXScvRbYh9X9okLuX9YMnz1HR4WgRTU2hocjBs15nmCNG/go-libp2p-floodsub"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core"
	"io"
)

func Subscribe(node *core.IpfsNode, topic string, didChange func(peer peer.ID, data []byte)) (*floodsub.Subscription, error) {
	if node.Floodsub == nil {
		return nil, errors.New("pub sub must be enabled.")
	}
	sub, err := node.Floodsub.Subscribe(topic)
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			msg, err := sub.Next(node.Context())
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
	if node.Floodsub == nil {
		return errors.New("pub sub must be enabled.")
	}
	logger.Debugf("Publishing: %s", topic)
	return node.Floodsub.Publish(topic, data)
}

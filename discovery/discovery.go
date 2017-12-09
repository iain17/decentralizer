package discovery

import (
	"context"
	"github.com/iain17/decentralizer/network"
	"time"
	"net"
	"github.com/iain17/timeout"
)

type Discovery struct {
	max int//Once we've reached it we won't engage new connections. ones connecting to us will trigger dropping the oldest connection.

	ctx context.Context
	cancel context.CancelFunc
	network *network.Network
	LocalNode *LocalNode
}

func New(network *network.Network, max int) (*Discovery, error) {
	ctx, cancel := context.WithCancel(context.Background())
	self := &Discovery{
		max: max,
		ctx: ctx,
		cancel: cancel,
		network: network,
	}
	var err error
	self.LocalNode, err = newLocalNode(self)
	if err != nil {
		return nil, err
	}
	return self, nil
}

func (d *Discovery) Stop() {
	d.cancel()
}

func (d *Discovery) WaitForPeers(num int, expire time.Duration) []*RemoteNode {
	timeout.Do(func(ctx context.Context) {
		for d.LocalNode.netTableService.peers.Len() < num {
			time.Sleep(100 * time.Millisecond)
		}
	}, expire)
	return d.LocalNode.netTableService.GetPeers()
}

func (d *Discovery) GetIP() net.IP {
	for {
		if d.LocalNode.ip != "" {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	return net.ParseIP(d.LocalNode.ip)
}

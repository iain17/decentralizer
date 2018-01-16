package discovery

import (
	"context"
	"github.com/iain17/discovery/network"
	"time"
	"net"
	"github.com/iain17/timeout"
	"github.com/iain17/logger"
)

type Discovery struct {
	max int//Once we've reached it we won't engage new connections. ones connecting to us will trigger dropping the oldest connection.

	ctx context.Context
	cancel context.CancelFunc
	network *network.Network
	LocalNode *LocalNode
	limited bool//Means we are on a limited connection. Means we won't advertise on DHT
}

func New(ctx context.Context, network *network.Network, max int, limitedConnection bool) (*Discovery, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok := r.(error)
			if !ok {
				logger.Errorf("panic in discovery package: %v", err)
			}
		}
	}()
	self := &Discovery{
		max: max,
		ctx: ctx,
		cancel: cancel,
		network: network,
		limited: limitedConnection,
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
	d.LocalNode.netTableService.Stop()
}

func (d *Discovery) WaitForPeers(num int, expire time.Duration) []*RemoteNode {
	timeout.Do(func(ctx context.Context) {
		d.LocalNode.waitTilReady()
		for d.LocalNode.netTableService.peers.Len() < num {
			time.Sleep(100 * time.Millisecond)
		}
	}, expire)
	return d.LocalNode.netTableService.GetPeers()
}

func (d *Discovery) GetIP() net.IP {
	d.LocalNode.waitTilReady()
	return net.ParseIP(d.LocalNode.ip)
}

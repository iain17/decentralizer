package discovery

import (
	"context"
	"github.com/iain17/discovery/network"
	"time"
	"net"
	"github.com/iain17/timeout"
	"github.com/iain17/logger"
)

type discoveredCB func(peer *RemoteNode)

type Discovery struct {
	max int//Once we've reached it we won't engage new connections. ones connecting to us will trigger dropping the oldest connection.

	ctx context.Context
	cancel context.CancelFunc
	network *network.Network
	LocalNode *LocalNode
	limited bool//Means we are on a limited connection. Means we won't advertise on DHT
	PeerDiscovered discoveredCB
}

func New(ctx context.Context, network *network.Network, max int, cb discoveredCB, limitedConnection bool, info map[string]string) (*Discovery, error) {
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
	if limitedConnection {
		CONCCURENT_NEW_CONNECTION = CONCCURENT_NEW_CONNECTION_LIMITED
	}
	self := &Discovery{
		max: max,
		ctx: ctx,
		cancel: cancel,
		network: network,
		limited: limitedConnection,
		PeerDiscovered: cb,
	}
	var err error
	self.LocalNode, err = newLocalNode(self)
	if err != nil {
		return nil, err
	}
	self.LocalNode.info = info
	return self, nil
}

//Make sure you call this before you quit, Just signalling the context won't be enough.
func (d *Discovery) Stop() {
	d.LocalNode.netTableService.Stop()
	d.cancel()
}

func (d *Discovery) WaitForPeers(num int, expire time.Duration) []*RemoteNode {
	d.LocalNode.waitTilReady()
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

func (d *Discovery) SetNetworkMessage(message string) {
	d.LocalNode.discoveryIRC.message = message
}

func (d *Discovery) GetNetworkMessages() []string {
	var messages []string
	for _, key := range d.LocalNode.discoveryIRC.messages.Keys() {
		if message, ok := d.LocalNode.discoveryIRC.messages.Get(key); ok {
			messages = append(messages, message.(string))
		}
	}
	return messages
}

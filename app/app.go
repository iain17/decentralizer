package app

import (
	"github.com/iain17/decentralizer/network"
	"github.com/iain17/decentralizer/discovery"
	"github.com/ipfs/go-ipfs/core"
	"fmt"
	"time"
)

type Decentralizer struct {
	n *network.Network
	d *discovery.Discovery
	i *core.IpfsNode
}

func New(networkStr string) (*Decentralizer, error) {
	n, err := network.UnmarshalFromPrivateKey(networkStr)
	if err != nil {
		return nil, err
	}
	d, err := discovery.New(n, MAX_DISCOVERED_PEERS)
	if err != nil {
		return nil, err
	}
	//Demo purposes
	path := fmt.Sprintf("/tmp/ipfs/%d", time.Now().Unix())
	i, err := OpenIPFSRepo(path, -1)
	if err != nil {
		return nil, err
	}
	instance := &Decentralizer{
		n: n,
		d: d,
		i: i,
	}
	instance.i.Bootstrap(core.BootstrapConfig{
		MinPeerThreshold:  4,
		Period:            30 * time.Second,
		ConnectionTimeout: (30 * time.Second) / 3, // Period / 3
		BootstrapPeers: instance.bootstrap,
	})
	return instance, nil
}
package app

import (
	"github.com/iain17/decentralizer/network"
	"github.com/iain17/decentralizer/discovery"
	"github.com/ipfs/go-ipfs/core"
	"fmt"
	"time"
	"github.com/iain17/decentralizer/app/ipfs"
	"github.com/iain17/decentralizer/app/pb"
	"github.com/iain17/logger"
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
	i, err := ipfs.OpenIPFSRepo(path, -1)
	if err != nil {
		return nil, err
	}
	instance := &Decentralizer{
		n: n,
		d: d,
		i: i,
	}
	logger.Infof("Our DiD is: %d", pb.GetPeer(i.Identity).DId)
	instance.i.Bootstrap(core.BootstrapConfig{
		MinPeerThreshold:  4,
		Period:            30 * time.Second,
		ConnectionTimeout: (30 * time.Second) / 3, // Period / 3
		BootstrapPeers: instance.bootstrap,
	})
	return instance, nil
}

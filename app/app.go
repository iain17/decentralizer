package app

import (
	"github.com/iain17/decentralizer/network"
	"github.com/iain17/decentralizer/discovery"
	"github.com/ipfs/go-ipfs/core"
	"time"
	"github.com/iain17/decentralizer/app/ipfs"
	"github.com/iain17/decentralizer/app/pb"
	"github.com/iain17/logger"
	"github.com/shibukawa/configdir"
	"gx/ipfs/QmTm7GoSkSSQPP32bZhvu17oY1AfvPKND6ELUdYAcKuR1j/floodsub"
	"errors"
)

type Decentralizer struct {
	n *network.Network
	d *discovery.Discovery
	i *core.IpfsNode

	sessions	  map[uint64]*pb.SessionInfo
	subscriptions map[uint32]*floodsub.Subscription
}

var configPath = configdir.New("ECorp", "Decentralizer")

func New(networkStr string) (*Decentralizer, error) {
	n, err := network.UnmarshalFromPrivateKey(networkStr)
	if err != nil {
		return nil, err
	}
	d, err := discovery.New(n, MAX_DISCOVERED_PEERS)
	if err != nil {
		return nil, err
	}
	paths := configPath.QueryFolders(configdir.System)
	if len(paths) == 0 {
		return nil, errors.New("queryFolder request failed")
	}
	i, err := ipfs.OpenIPFSRepo(paths[0].Path, -1)
	if err != nil {
		return nil, err
	}
	instance := &Decentralizer{
		n: n,
		d: d,
		i: i,
	}
	_, dID := pb.GetPeer(i.Identity)
	logger.Infof("Our DiD is: %v", dID)
	instance.i.Bootstrap(core.BootstrapConfig{
		MinPeerThreshold:  4,
		Period:            30 * time.Second,
		ConnectionTimeout: (30 * time.Second) / 3, // Period / 3
		BootstrapPeers: instance.bootstrap,
	})
	return instance, nil
}

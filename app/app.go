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
	"github.com/iain17/decentralizer/app/sessionstore"
	"gx/ipfs/QmTm7GoSkSSQPP32bZhvu17oY1AfvPKND6ELUdYAcKuR1j/floodsub"
)

type Decentralizer struct {
	n *network.Network
	d *discovery.Discovery
	i *core.IpfsNode

	sessions map[uint32]*sessionstore.Store
	subscriptions map[uint32]*floodsub.Subscription
}

var configPath = configdir.New("eCORp", "Decentralizer")

func New(networkStr string) (*Decentralizer, error) {
	n, err := network.UnmarshalFromPrivateKey(networkStr)
	if err != nil {
		return nil, err
	}
	d, err := discovery.New(n, MAX_DISCOVERED_PEERS)
	if err != nil {
		return nil, err
	}
	i, err := ipfs.OpenIPFSRepo(configPath.LocalPath, -1)
	if err != nil {
		return nil, err
	}
	instance := &Decentralizer{
		n: n,
		d: d,
		i: i,
		sessions: make(map[uint32]*sessionstore.Store),
		subscriptions: make(map[uint32]*floodsub.Subscription),
	}
	logger.Infof("Our DiD is: %v", pb.GetPeer(i.Identity))
	instance.i.Bootstrap(core.BootstrapConfig{
		MinPeerThreshold:  4,
		Period:            30 * time.Second,
		ConnectionTimeout: (30 * time.Second) / 3, // Period / 3
		BootstrapPeers: instance.bootstrap,
	})
	go instance.Advertise()
	return instance, nil
}

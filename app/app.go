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
)

type Decentralizer struct {
	n *network.Network
	d *discovery.Discovery
	i *core.IpfsNode

	sessions *sessionstore.Store
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
	sessions, err := sessionstore.New(1000, time.Duration((EXPIRE_TIME_SESSION * 1.5) * time.Second))
	if err != nil {
		return nil, err
	}
	instance := &Decentralizer{
		n: n,
		d: d,
		i: i,
		sessions: sessions,
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

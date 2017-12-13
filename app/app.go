package app

import (
	"github.com/iain17/decentralizer/app/ipfs"
	"github.com/iain17/decentralizer/discovery"
	"github.com/iain17/decentralizer/network"
	"github.com/iain17/logger"
	"github.com/shibukawa/configdir"
	"gx/ipfs/QmNUKMfTHQQpEwE8bUdv5qmKC3ymdW7zw82LFS8D6MQXmu/go-ipfs/core"
	"time"
	//"gx/ipfs/QmTm7GoSkSSQPP32bZhvu17oY1AfvPKND6ELUdYAcKuR1j/floodsub"
	"errors"
	"github.com/iain17/decentralizer/app/sessionstore"
	"net"
)

type Decentralizer struct {
	n *network.Network
	d *discovery.Discovery
	i *core.IpfsNode
	b *ipfs.BitswapService

	sessions               map[uint64]*sessionstore.Store
	sessionIdToSessionType map[uint64]uint64
	//subscriptions map[uint32]*floodsub.Subscription
}

var configPath = configdir.New("ECorp", "Decentralizer")

func getIpfsPath() (string, error) {
	paths := configPath.QueryFolders(configdir.Global)
	if len(paths) == 0 {
		return "", errors.New("queryFolder request failed")
	}
	return paths[0].Path, nil
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
	path, err := getIpfsPath()
	if err != nil {
		return nil, err
	}
	i, err := ipfs.OpenIPFSRepo(path, -1)
	if err != nil {
		return nil, err
	}
	b, err := ipfs.NewBitSwap(i)
	if err != nil {
		return nil, err
	}
	instance := &Decentralizer{
		n:                      n,
		d:                      d,
		i:                      i,
		b:                      b,
		sessions:               make(map[uint64]*sessionstore.Store),
		sessionIdToSessionType: make(map[uint64]uint64),
	}
	instance.initMatchmaking()
	instance.initMessaging()
	_, dnID := PeerToDnId(i.Identity)
	logger.Infof("Our dnID is: %v", dnID)
	go instance.i.Bootstrap(core.BootstrapConfig{
		MinPeerThreshold:  4,
		Period:            30 * time.Second,
		ConnectionTimeout: (30 * time.Second) / 3, // Period / 3
		BootstrapPeers:    instance.bootstrap,
	})
	return instance, nil
}

func (d *Decentralizer) GetIP() net.IP {
	if d.d != nil {
		return d.d.GetIP()
	}
	return net.ParseIP("127.0.0.1")
}

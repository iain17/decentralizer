package app

import (
	"github.com/iain17/decentralizer/app/ipfs"
	//"github.com/iain17/discovery"
	"github.com/iain17/discovery/network"
	"github.com/iain17/logger"
	"github.com/shibukawa/configdir"
	"gx/ipfs/QmNUKMfTHQQpEwE8bUdv5qmKC3ymdW7zw82LFS8D6MQXmu/go-ipfs/core"
	"time"
	"errors"
	"github.com/iain17/decentralizer/app/sessionstore"
	"net"
	"github.com/iain17/decentralizer/app/peerstore"
	"context"
	"github.com/ccding/go-stun/stun"
)

type Decentralizer struct {
	n *network.Network
	//d *discovery.Discovery
	i *core.IpfsNode
	b *ipfs.BitswapService

	ip 					   *net.IP
	sessions               map[uint64]*sessionstore.Store
	sessionIdToSessionType map[uint64]uint64
	peers			   	   *peerstore.Store
	directMessage		   chan *DirectMessage
}

var configPath = configdir.New("ECorp", "Decentralizer")

func getIpfsPath() (string, error) {
	paths := configPath.QueryFolders(configdir.Global)
	if len(paths) == 0 {
		return "", errors.New("queryFolder request failed")
	}
	return paths[0].Path, nil
}

func New(ctx context.Context, networkStr string, privateKey bool) (*Decentralizer, error) {
	var n *network.Network
	var err error
	if privateKey {
		n, err = network.UnmarshalFromPrivateKey(networkStr)
	} else {
		n, err = network.Unmarshal(networkStr)
	}
	if err != nil {
		return nil, err
	}
	//d, err := discovery.New(n, MAX_DISCOVERED_PEERS)
	//if err != nil {
	//	return nil, err
	//}
	path, err := getIpfsPath()
	if err != nil {
		return nil, err
	}
	i, err := ipfs.OpenIPFSRepo(ctx, path, -1)
	if err != nil {
		return nil, err
	}
	b, err := ipfs.NewBitSwap(i)
	if err != nil {
		return nil, err
	}
	peers, err := peerstore.New(MAX_CONTACTS, time.Duration((EXPIRE_TIME_CONTACT*1.5)*time.Second), i.Identity)
	if err != nil {
		return nil, err
	}
	instance := &Decentralizer{
		n:                      n,
		//d:                      d,
		i:                      i,
		b:                      b,
		sessions:               make(map[uint64]*sessionstore.Store),
		sessionIdToSessionType: make(map[uint64]uint64),
		peers:				    peers,
		directMessage: 			make(chan *DirectMessage, 10),
	}
	instance.initMatchmaking()
	instance.initMessaging()
	instance.initAddressbook()
	_, dnID := peerstore.PeerToDnId(i.Identity)
	logger.Infof("Our dnID is: %v", dnID)
	//go instance.i.Bootstrap(core.BootstrapConfig{
	//	MinPeerThreshold:  4,
	//	Period:            30 * time.Second,
	//	ConnectionTimeout: (30 * time.Second) / 3, // Period / 3
	//	BootstrapPeers:    instance.bootstrap,
	//})
	return instance, nil
}

func (d *Decentralizer) GetIP() net.IP {
	if d.ip == nil {
		client := stun.NewClient()
		_, host, _ := client.Discover()
		ip := net.ParseIP(host.IP())
		d.ip = &ip
	}
	return *d.ip
}

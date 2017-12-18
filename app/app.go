package app

import (
	"github.com/iain17/decentralizer/app/ipfs"
	//"github.com/iain17/discovery"
	"context"
	"errors"
	"github.com/ccding/go-stun/stun"
	"github.com/iain17/decentralizer/app/peerstore"
	"github.com/iain17/decentralizer/app/sessionstore"
	"github.com/iain17/discovery/network"
	"github.com/iain17/logger"
	"github.com/shibukawa/configdir"
	"gx/ipfs/QmTxUjSZnG7WmebrX2U7furEPNSy33pLgA53PtpJYJSZSn/go-ipfs/core"
	"net"
	"time"
)

type Decentralizer struct {
	n *network.Network
	//d *discovery.Discovery
	i                      *core.IpfsNode
	b                      *ipfs.BitswapService
	ip                     net.IP
	sessions               map[uint64]*sessionstore.Store
	sessionIdToSessionType map[uint64]uint64
	peers                  *peerstore.Store
	directMessage          chan *DirectMessage
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
	client := stun.NewClient()
	_, host, err := client.Discover()
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
		n:   n,
		//d:                      d,
		i:                      i,
		b:                      b,
		ip:                     net.ParseIP(host.IP()),
		sessions:               make(map[uint64]*sessionstore.Store),
		sessionIdToSessionType: make(map[uint64]uint64),
		peers:         peers,
		directMessage: make(chan *DirectMessage, 10),
	}
	instance.initMatchmaking()
	instance.initMessaging()
	instance.initAddressbook()
	_, dnID := peerstore.PeerToDnId(i.Identity)
	logger.Infof("Our dnID is: %v", dnID)
	bs := core.DefaultBootstrapConfig
	bs.BootstrapPeers = nil //instance.bootstrap
	instance.i.Bootstrap(bs)
	return instance, nil
}

func (d *Decentralizer) GetIP() net.IP {
	return d.ip
}

func (s *Decentralizer) Stop() {
	if s.i != nil {
		s.i.Close()
	}
}

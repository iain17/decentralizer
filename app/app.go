package app

import (
	"github.com/iain17/decentralizer/app/ipfs"
	//"github.com/iain17/discovery"
	"context"
	"errors"
	"github.com/iain17/decentralizer/app/peerstore"
	"github.com/iain17/decentralizer/app/sessionstore"
	"github.com/iain17/discovery/network"
	"github.com/iain17/logger"
	"github.com/shibukawa/configdir"
	"gx/ipfs/QmTxUjSZnG7WmebrX2U7furEPNSy33pLgA53PtpJYJSZSn/go-ipfs/core"
	libp2pPeer "gx/ipfs/QmWNY7dV54ZDYmTA1ykVdwNCqC11mpU4zSUp6XDpLTH9eG/go-libp2p-peer"
	"net"
	"time"
	"github.com/ccding/go-stun/stun"
	"github.com/robfig/cron"
	"github.com/iain17/discovery"
	"github.com/iain17/decentralizer/pb"
	"os"
)

type Decentralizer struct {
	ctx 				   context.Context
	n 					   *network.Network
	cron				   *cron.Cron
	d					   *discovery.Discovery
	i                      *core.IpfsNode
	b                      *ipfs.BitswapService
	ip                     *net.IP

	//Peer ids that did not respond to our queries.
	ignore 				   map[string]bool

	//Matchmaking
	sessions               map[uint64]*sessionstore.Store
	sessionIdToSessionType map[uint64]uint64
	searches 			   map[uint64]*search

	//addressbook
	peers                  *peerstore.Store
	addressBookChanged     bool

	//messaging
	DirectMessage          chan *pb.RPCDirectMessage

	//Publisher files
	publisherUpdate  	   *pb.PublisherUpdate
	publisherDefinition	   *pb.PublisherDefinition
	searchingForPublisherUpdate bool
}

var configPath = configdir.New("ECorp", "Decentralizer")

func getIpfsPath() string {
	paths := configPath.QueryFolders(configdir.Global)
	if len(paths) == 0 {
		panic(errors.New("queryFolder request failed"))
	}
	return paths[0].Path
}

func Reset() {
	os.RemoveAll(configPath.QueryCacheFolder().Path)
	os.RemoveAll(getIpfsPath())
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
	var d *discovery.Discovery
	if USE_OWN_BOOTSTRAPPING {
		d, err = discovery.New(n, MAX_DISCOVERED_PEERS)
		if err != nil {
			return nil, err
		}
	}
	ipfsPath := getIpfsPath()
	logger.Infof("IPFS path: %s", ipfsPath)
	logger.Infof("Cache path: %s", configPath.QueryCacheFolder().Path)
	i, err := ipfs.OpenIPFSRepo(ctx, ipfsPath, -1)
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
		ctx:					ctx,
		cron: 				   cron.New(),
		n:   					n,
		d:                      d,
		i:                      i,
		b:                      b,
		sessions:               make(map[uint64]*sessionstore.Store),
		sessionIdToSessionType: make(map[uint64]uint64),
		searches:				make(map[uint64]*search),
		peers:         			peers,
		DirectMessage: 			make(chan *pb.RPCDirectMessage, 10),
		ignore:					make(map[string]bool),
	}
	err = instance.bootstrap()
	if err == nil {
		reveries, _ := Asset("reveries.flac")
		instance.SavePeerFile("reveries.flac", reveries)
	}

	instance.GetIP()
	instance.initMatchmaking()
	instance.initMessaging()
	instance.initAddressbook()
	instance.initPublisherFiles()
	instance.cron.Start()
	_, dnID := peerstore.PeerToDnId(i.Identity)
	logger.Infof("Our dnID is: %v", dnID)
	return instance, err
}

func (s *Decentralizer) decodePeerId(id string) (libp2pPeer.ID, error) {
	if id == "self" {
		return s.i.Identity, nil
	}
	return libp2pPeer.IDB58Decode(id)
}

func (d *Decentralizer) GetIP() net.IP {
	if d.d != nil {
		return d.d.GetIP()
	}
	if d.ip == nil {
		stun := stun.NewClient()
		nat, host, err := stun.Discover()
		if err != nil {
			logger.Error(err)
			time.Sleep(5 * time.Second)
			return d.GetIP()
		}
		logger.Infof("NAT type: %s", nat.String())
		ip := net.ParseIP(host.IP())
		d.ip = &ip
	}

	return *d.ip
}

func (s *Decentralizer) Stop() {
	s.cron.Stop()
	if s.i != nil {
		s.i.Close()
	}
}

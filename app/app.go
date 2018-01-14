package app

import (
	"github.com/iain17/decentralizer/app/ipfs"
	"context"
	"errors"
	"github.com/iain17/decentralizer/app/peerstore"
	"github.com/iain17/decentralizer/app/sessionstore"
	"github.com/iain17/discovery/network"
	"github.com/iain17/logger"
	"github.com/shibukawa/configdir"
	"gx/ipfs/QmYHpXQEWuhwgRFBnrf4Ua6AZhcqXCYa7Biv65SLGgTgq5/go-ipfs/core"
	libp2pPeer "gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
	"net"
	"time"
	"github.com/ccding/go-stun/stun"
	"github.com/jasonlvhit/gocron"
	"github.com/iain17/discovery"
	"github.com/iain17/decentralizer/pb"
	"os"
	coreiface "gx/ipfs/QmYHpXQEWuhwgRFBnrf4Ua6AZhcqXCYa7Biv65SLGgTgq5/go-ipfs/core/coreapi/interface"
	"sync"
	"gx/ipfs/QmYHpXQEWuhwgRFBnrf4Ua6AZhcqXCYa7Biv65SLGgTgq5/go-ipfs/core/coreapi"
	"github.com/iain17/kvcache/lttlru"
	"github.com/spf13/afero"
	"github.com/hashicorp/golang-lru"
	"hash"
	"hash/crc32"
)

type Decentralizer struct {
	mutex				   sync.Mutex
	ctx 				   context.Context
	n 					   *network.Network
	cron				   *gocron.Scheduler
	cronChan			   chan bool
	d					   *discovery.Discovery
	i                      *core.IpfsNode
	b                      *ipfs.BitswapService
	ip                     *net.IP
	api 				   coreiface.CoreAPI

	//Peer ids that did not respond to our queries.
	ignore 				   *lttlru.LruWithTTL
	crcTable			   hash.Hash32
	unmarshalCache		   *lru.Cache//We unmarshal the same data over and over. Let us cache this.

	//Storage
	filesApi       		   *ipfs.FilesAPI
	ufs					   afero.Fs

	//Matchmaking
	matchmakingMutex	   sync.Mutex
	searchMutex			   sync.Mutex
	sessionQueries		   chan sessionRequest
	sessions               map[uint64]*sessionstore.Store
	sessionIdToSessionType map[uint64]uint64
	searches 			   *lttlru.LruWithTTL

	//addressbook
	peers                  *peerstore.Store
	addressBookChanged     bool

	//messaging
	directMessageChannels  map[uint32]chan *pb.RPCDirectMessage

	//Publisher files
	publisherRecord     *pb.DNPublisherRecord
	publisherDefinition *pb.PublisherDefinition
}

var configPath = configdir.New("ECorp", "Decentralizer")
var Base = getBasePath()

func getBasePath() *configdir.Config {
	paths := configPath.QueryFolders(configdir.Global)
	if len(paths) == 0 {
		panic(errors.New("queryFolder request failed"))
	}
	return paths[0]
}

func Reset() {
	os.RemoveAll(configPath.QueryCacheFolder().Path)
	os.RemoveAll(getBasePath().Path+"/ipfs")
}

func New(ctx context.Context, networkStr string, privateKey bool, limitedConnection bool) (*Decentralizer, error) {
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
	ipfsPath := Base.Path+"/ipfs"
	logger.Infof("IPFS path: %s", ipfsPath)
	logger.Infof("Cache path: %s", configPath.QueryCacheFolder().Path)
	i, err := ipfs.OpenIPFSRepo(ctx, ipfsPath, limitedConnection)
	if err != nil {
		return nil, err
	}
	var d *discovery.Discovery
	if USE_OWN_BOOTSTRAPPING {
		d, err = discovery.New(ctx, n, MAX_DISCOVERED_PEERS, limitedConnection)
		if err != nil {
			return nil, err
		}
	}
	b, err := ipfs.NewBitSwap(i)
	if err != nil {
		return nil, err
	}
	ignore, err := lttlru.NewTTL(MAX_IGNORE)
	if err != nil {
		return nil, err
	}
	unmarshalCache, err := lru.New(MAX_UNMARSHAL_CACHE)
	if err != nil {
		return nil, err
	}
	instance := &Decentralizer{
		ctx:					ctx,
		cron: 				    gocron.NewScheduler(),
		n:   					n,
		d:                      d,
		i:                      i,
		b:                      b,
		api:					coreapi.NewCoreAPI(i),
		directMessageChannels:  make(map[uint32]chan *pb.RPCDirectMessage),
		ignore:					ignore,
		unmarshalCache:			unmarshalCache,
		crcTable:				crc32.NewIEEE(),
	}
	go instance.bootstrap()

	instance.initStorage()
	instance.initMatchmaking()
	instance.initMessaging()
	instance.initAddressbook()
	instance.initPublisherFiles()
	instance.cronChan = instance.cron.Start()
	return instance, err
}

func (s *Decentralizer) decodePeerId(id string) (libp2pPeer.ID, error) {
	if id == "self" {
		return s.i.Identity, nil
	}
	return libp2pPeer.IDB58Decode(id)
}

func (d *Decentralizer) GetIP() net.IP {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if d.d != nil {
		ip := d.d.GetIP()
		d.ip = &ip
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

func (d *Decentralizer) IsEnoughPeers() bool {
	lenPeers := len(d.i.PeerHost.Network().Peers())
	return lenPeers >= MIN_CONNECTED_PEERS
}

func (d *Decentralizer) WaitTilEnoughPeers() {
	for {
		if d.IsEnoughPeers() {
			break
		}
		time.Sleep(300 * time.Millisecond)
	}
}

func (s *Decentralizer) Stop() {
	s.cronChan <- false
	if s.i != nil {
		s.i.Close()
	}
}

package app

import (
	"context"
	"github.com/hashicorp/golang-lru"
	"github.com/iain17/decentralizer/app/ipfs"
	"github.com/iain17/decentralizer/pb"
	"github.com/iain17/discovery/network"
	"github.com/iain17/kvcache/lttlru"
	"github.com/iain17/logger"
	"github.com/jasonlvhit/gocron"
	"github.com/shibukawa/configdir"
	"github.com/spf13/afero"
	logging "gx/ipfs/QmcVVHfdyv15GVPk7NrxdWjh2hLVccXnoD8j2tyQShiXJb/go-log"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core/coreapi"
	"hash/crc32"
	"net"
	"os"
	"time"
)

var testNetwork *network.Network
var testSlaveNetwork *network.Network //just the public key

func init() {
	MIN_CONNECTED_PEERS = 1
	logger.AddOutput(logger.Stdout{
		MinLevel: logger.INFO, //logger.DEBUG,
		Colored:  true,
	})
	logging.Configure(logging.LevelInfo)
	configPath = configdir.New("ECorp", "Decentralizer-test")
	os.RemoveAll(configPath.QueryCacheFolder().Path)
	testNetwork, _ = network.New()
	testSlaveNetwork, _ = network.Unmarshal(testNetwork.Marshal())
	Base = getBasePath()
}

func fakeNew(ctx context.Context, node *core.IpfsNode, master bool) *Decentralizer {
	os.RemoveAll(configPath.QueryCacheFolder().Path)
	b, err := ipfs.NewBitSwap(node)
	if err != nil {
		panic(err)
	}

	//Build a new network.
	var n *network.Network
	if master {
		n = testNetwork
	} else {
		n = testSlaveNetwork
	}
	ignore, err := lttlru.NewTTL(MAX_IGNORE)
	if err != nil {
		panic(err)
	}
	unmarshalCache, err := lru.New(MAX_UNMARSHAL_CACHE)
	if err != nil {
		panic(err)
	}
	ip := net.ParseIP("127.0.0.1")
	instance := &Decentralizer{
		ctx:  ctx,
		cron: gocron.NewScheduler(),
		n:    n,
		ip:   &ip,
		i:    node,
		b:    b,
		api:  coreapi.NewCoreAPI(node),
		directMessageChannels: make(map[uint32]chan *pb.RPCDirectMessage),
		ignore:                ignore,
		unmarshalCache:        unmarshalCache,
		crcTable:              crc32.NewIEEE(),
	}
	//Mock filesystem
	instance.peerFileSystem = afero.NewMemMapFs()
	instance.fs = instance.peerFileSystem
	instance.WaitTilEnoughPeers()
	Base = &configdir.Config{
		Type: configdir.Cache,
		Path: "/tmp/" + time.Now().Format("20060102150405"),
	}

	instance.initializeComponents(true)

	go func() {
		<-instance.ctx.Done()
		instance.cronChan <- false
	}()

	return instance
}

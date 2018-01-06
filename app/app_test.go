package app

import (
	"github.com/iain17/decentralizer/app/ipfs"
	"github.com/iain17/logger"
	"gx/ipfs/QmYHpXQEWuhwgRFBnrf4Ua6AZhcqXCYa7Biv65SLGgTgq5/go-ipfs/core"
	"net"
	"github.com/jasonlvhit/gocron"
	"github.com/shibukawa/configdir"
	"os"
	"github.com/iain17/discovery/network"
	"github.com/iain17/decentralizer/pb"
	"github.com/iain17/kvcache/lttlru"
	"gx/ipfs/QmYHpXQEWuhwgRFBnrf4Ua6AZhcqXCYa7Biv65SLGgTgq5/go-ipfs/core/coreapi"
	"github.com/spf13/afero"
	logging "gx/ipfs/QmSpJByNKFX1sCsHBEp3R73FL4NF6FnQTEGyNAXHm2GS52/go-log"
	"context"
)

var testNetwork *network.Network
var testSlaveNetwork *network.Network//just the public key

func init() {
	MIN_CONNECTED_PEERS = 1
	logger.AddOutput(logger.Stdout{
		MinLevel: logger.DEBUG, //logger.DEBUG,
		Colored:  true,
	})
	logging.Configure(logging.LevelDebug)
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

	ip := net.ParseIP("127.0.0.1")
	instance := &Decentralizer{
		ctx:					ctx,
		cron:					gocron.NewScheduler(),
		n:						n,
		ip:                     &ip,
		i:                      node,
		b:                      b,
		api:					coreapi.NewCoreAPI(node),
		directMessageChannels: 	make(map[uint32]chan *pb.RPCDirectMessage),
		ignore:					ignore,
	}
	instance.cronChan = instance.cron.Start()
	instance.initStorage()
	instance.initMatchmaking()
	instance.initMessaging()
	instance.initAddressbook()
	instance.initPublisherFiles()

	go func() {
		<- instance.ctx.Done()
		instance.cronChan <- false
	}()

	//Mock UFS
	instance.ufs = afero.NewMemMapFs()
	instance.WaitTilEnoughPeers()

	return instance
}

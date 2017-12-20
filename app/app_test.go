package app

import (
	"github.com/iain17/decentralizer/app/ipfs"
	"github.com/iain17/decentralizer/app/sessionstore"
	"github.com/iain17/logger"
	"gx/ipfs/QmTxUjSZnG7WmebrX2U7furEPNSy33pLgA53PtpJYJSZSn/go-ipfs/core"
	//logging "gx/ipfs/QmSpJByNKFX1sCsHBEp3R73FL4NF6FnQTEGyNAXHm2GS52/go-log"
	"github.com/iain17/decentralizer/app/peerstore"
	"net"
	"github.com/robfig/cron"
	"time"
	"github.com/shibukawa/configdir"
	"os"
)

func init() {
	logger.AddOutput(logger.Stdout{
		MinLevel: logger.INFO, //logger.DEBUG,
		Colored:  true,
	})
	//logging.Configure(logging.LevelDebug)
	configPath = configdir.New("ECorp", "Decentralizer-test")
	os.RemoveAll(configPath.QueryCacheFolder().Path)
}

func fakeNew(node *core.IpfsNode) *Decentralizer {
	b, err := ipfs.NewBitSwap(node)
	if err != nil {
		panic(err)
	}
	peers, err := peerstore.New(MAX_CONTACTS, time.Duration((EXPIRE_TIME_CONTACT*1.5)*time.Second), node.Identity)
	if err != nil {
		panic(err)
	}
	ip := net.ParseIP("127.0.0.1")
	instance := &Decentralizer{
		cron:					cron.New(),
		ip:                     &ip,
		i:                      node,
		b:                      b,
		sessions:               make(map[uint64]*sessionstore.Store),
		sessionIdToSessionType: make(map[uint64]uint64),
		searches:				make(map[uint64]*search),
		peers:         peers,
		directMessage: make(chan *DirectMessage, 10),
	}
	instance.initMatchmaking()
	instance.initMessaging()
	instance.initAddressbook()
	return instance
}

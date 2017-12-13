package app

import (
	"github.com/ipfs/go-ipfs/core"
	"github.com/iain17/decentralizer/app/sessionstore"
	"github.com/iain17/decentralizer/app/ipfs"
	"github.com/iain17/logger"
	logging "gx/ipfs/QmSpJByNKFX1sCsHBEp3R73FL4NF6FnQTEGyNAXHm2GS52/go-log"
)

func init() {
	logger.AddOutput(logger.Stdout{
		MinLevel: logger.INFO, //logger.DEBUG,
		Colored:  true,
	})
	logging.Configure(logging.LevelDebug)
}

func fakeNew(node *core.IpfsNode) *Decentralizer {
	b, err := ipfs.NewBitSwap(node)
	if err != nil {
		return nil
	}
	instance := &Decentralizer{
		i: node,
		b: b,
		sessions: make(map[uint64]*sessionstore.Store),
		sessionIdToSessionType: make(map[uint64]uint64),
	}
	instance.initMatchmaking()
	instance.initMessaging()
	return instance
}

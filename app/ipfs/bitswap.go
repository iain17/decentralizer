package ipfs

import (
	"github.com/giantswarm/retry-go"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/exchange/bitswap"
	"gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
	"reflect"
	bsnet "github.com/ipfs/go-ipfs/exchange/bitswap/network"
	"gx/ipfs/QmNp85zy9RLrQ5oQD4hPyS39ezrrXpcaa7R4Y9kxdWQLLQ/go-cid"
	"unsafe"
	"gx/ipfs/QmSn9Td7xgxm9EV7iEjTckpUWmWApggzPxu7eFGWkkpwin/go-block-format"
	logging "gx/ipfs/QmSpJByNKFX1sCsHBEp3R73FL4NF6FnQTEGyNAXHm2GS52/go-log"
	"time"
)

var log = logging.Logger("BitswapService")

//Find other peers around a subject.
//This is done by using the bitswap network of IPFS which is currently powered by DHT.

type BitswapService struct {
	node *core.IpfsNode
	network bsnet.BitSwapNetwork
}

func NewBitSwap(node *core.IpfsNode) (*BitswapService, error) {
	//Extract the network unexported value from the bitswap exchange of ipfs
	exchange := node.Exchange.(*bitswap.Bitswap)
	pointerVal := reflect.ValueOf(exchange)
	val := reflect.Indirect(pointerVal)

	member := val.FieldByName("network")
	ptrToY := unsafe.Pointer(member.UnsafeAddr())
	realPtrToY := (*bsnet.BitSwapNetwork)(ptrToY)
	network := *(realPtrToY)

	return &BitswapService{
		node: node,
		network: network,
	}, nil
}

func (b *BitswapService) Find(subject string, num int) <-chan peer.ID {
	log.Debugf("Find subject: %s", subject)
	return b.network.FindProvidersAsync(b.node.Context(), StringToCid(subject), num)
}

func (b *BitswapService) Provide(subject string) error {
	log.Debugf("Provide subject: %s", subject)
	return b.network.Provide(b.node.Context(), StringToCid(subject))
}

func StringToCid(value string) *cid.Cid {
	block := blocks.NewBlock([]byte(value))
	return block.Cid()
}
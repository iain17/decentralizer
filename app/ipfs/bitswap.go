package ipfs

import (
	"errors"
	"gx/ipfs/QmNUKMfTHQQpEwE8bUdv5qmKC3ymdW7zw82LFS8D6MQXmu/go-ipfs/core"
	"gx/ipfs/QmNUKMfTHQQpEwE8bUdv5qmKC3ymdW7zw82LFS8D6MQXmu/go-ipfs/exchange/bitswap"
	bsnet "gx/ipfs/QmNUKMfTHQQpEwE8bUdv5qmKC3ymdW7zw82LFS8D6MQXmu/go-ipfs/exchange/bitswap/network"
	"gx/ipfs/QmNp85zy9RLrQ5oQD4hPyS39ezrrXpcaa7R4Y9kxdWQLLQ/go-cid"
	logging "gx/ipfs/QmSpJByNKFX1sCsHBEp3R73FL4NF6FnQTEGyNAXHm2GS52/go-log"
	"gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
	"reflect"
	"unsafe"
	u "gx/ipfs/QmSU6eubNdhXjFBJBSksTp8kv8YRub8mGAPv8tVJHmL2EU/go-ipfs-util"
	"gx/ipfs/QmSn9Td7xgxm9EV7iEjTckpUWmWApggzPxu7eFGWkkpwin/go-block-format"
)

var log = logging.Logger("BitswapService")

//Find other peers around a subject.
//This is done by using the bitswap network of IPFS which is currently powered by DHT.

type BitswapService struct {
	node    *core.IpfsNode
	network bsnet.BitSwapNetwork
}

func NewBitSwap(node *core.IpfsNode) (*BitswapService, error) {
	//Extract the network unexported value from the bitswap exchange of ipfs
	if exchange, ok := node.Exchange.(*bitswap.Bitswap); ok {

		pointerVal := reflect.ValueOf(exchange)
		val := reflect.Indirect(pointerVal)

		member := val.FieldByName("network")
		ptrToY := unsafe.Pointer(member.UnsafeAddr())
		realPtrToY := (*bsnet.BitSwapNetwork)(ptrToY)
		network := *(realPtrToY)

		return &BitswapService{
			node:    node,
			network: network,
		}, nil
	} else {
		return nil, errors.New("interface conversion: node.Exchange is not *bitswap.Bitswap")
	}
}

func (b *BitswapService) Find(subject string, num int) <-chan peer.ID {
	log.Debugf("Find subject: %s", subject)
	peers := b.network.FindProvidersAsync(b.node.Context(), StringToCid2(subject), num)
	log.Debugf("Found %d around %s", len(peers), subject)
	return peers
}

func (b *BitswapService) Provide(subject string) error {
	log.Debugf("Provide subject: %s", subject)
	return b.network.Provide(b.node.Context(), StringToCid2(subject))
}

func StringToCid(value string) *cid.Cid {
	block := blocks.NewBlock([]byte(value))
	return block.Cid()
}

func StringToCid2(value string) *cid.Cid {
	return cid.NewCidV0(u.Hash([]byte(value)))
}

package ipfs

import (
	"errors"
	u "gx/ipfs/QmPsAfmDBnZN3kZGSuNwvCNDZiHneERSKmRcFyG3UkvcT3/go-ipfs-util"
	logging "gx/ipfs/QmSpJByNKFX1sCsHBEp3R73FL4NF6FnQTEGyNAXHm2GS52/go-log"
	"gx/ipfs/QmTxUjSZnG7WmebrX2U7furEPNSy33pLgA53PtpJYJSZSn/go-ipfs/core"
	"gx/ipfs/QmTxUjSZnG7WmebrX2U7furEPNSy33pLgA53PtpJYJSZSn/go-ipfs/exchange/bitswap"
	bsnet "gx/ipfs/QmTxUjSZnG7WmebrX2U7furEPNSy33pLgA53PtpJYJSZSn/go-ipfs/exchange/bitswap/network"
	"gx/ipfs/QmWNY7dV54ZDYmTA1ykVdwNCqC11mpU4zSUp6XDpLTH9eG/go-libp2p-peer"
	"gx/ipfs/QmYsEQydGrsxNZfAiskvQ76N2xE9hDQtSAkRSynwMiUK3c/go-block-format"
	"gx/ipfs/QmeSrf6pzut73u6zLQkRFQ3ygt3k6XFT2kjdYP8Tnkwwyg/go-cid"
	"reflect"
	"unsafe"
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

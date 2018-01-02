package ipfs

import (
	"errors"
	u "gx/ipfs/QmSU6eubNdhXjFBJBSksTp8kv8YRub8mGAPv8tVJHmL2EU/go-ipfs-util"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/exchange/bitswap"
	bsnet "github.com/ipfs/go-ipfs/exchange/bitswap/network"
	"gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
	"gx/ipfs/QmSn9Td7xgxm9EV7iEjTckpUWmWApggzPxu7eFGWkkpwin/go-block-format"
	"gx/ipfs/QmNp85zy9RLrQ5oQD4hPyS39ezrrXpcaa7R4Y9kxdWQLLQ/go-cid"
	"reflect"
	"unsafe"
	"github.com/iain17/timeout"
	"time"
	"context"
	"github.com/iain17/logger"
)

//Find other peers around a subject.
//This is done by using the bitswap network of IPFS which is currently powered by DHT.
//TODO: d.i.Routing?!?!
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
	logger.Debugf("Find subject: %s", subject)
	peers := b.network.FindProvidersAsync(b.node.Context(), StringToCid(subject), num)
	logger.Debugf("Found %d around %s", len(peers), subject)
	return peers
}

func (b *BitswapService) Provide(subject string) error {
	var err error
	completed := false
	timeout.Do(func(ctx context.Context) {
		err = b.network.Provide(b.node.Context(), StringToCid(subject))
		completed = true
		logger.Debugf("Provided subject: %s", subject)
	}, 5*time.Second)
	//if !completed {
	//	err = errors.New("could not provide '%s' in under 15 seconds. Check if you are connected to enough peers")
	//}
	return err
}

func StringToCid(value string) *cid.Cid {
	block := blocks.NewBlock([]byte(value))
	return block.Cid()
}

func StringToCid2(value string) *cid.Cid {
	return cid.NewCidV0(u.Hash([]byte(value)))
}

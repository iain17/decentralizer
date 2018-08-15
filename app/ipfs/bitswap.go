package ipfs

import (
	"errors"
	"fmt"
	"github.com/iain17/logger"
	u "gx/ipfs/QmPdKqUcHGFdeSpvjVoaTRPPstGif9GBZb5Q56RVw9o69A/go-ipfs-util"
	"gx/ipfs/QmVsp2KdPYE6M8ryzCk5KHLo3zprcY5hBDaYx6uPCFUdxA/go-libp2p-record"
	ipdht "gx/ipfs/QmTktQYCKzQjhxF6dk5xJPRuhHn3JBiKGvMLoiDy1mYmxC/go-libp2p-kad-dht"
	b58 "gx/ipfs/QmWFAMPqsEyUX7gDUsRVmMWz59FxSpJ1b2v6bJ1yYzo7jY/go-base58-fast/base58"
	pstore "gx/ipfs/QmZR2XWVVBCtbgBWnQhWk2xcQfaR3W8faQPriAiaaj7rsr/go-libp2p-peerstore"
	"gx/ipfs/QmYVNvtQkeZ6AKSwDrjQTs432QtL6umrrK41EBq3cu7iSP/go-cid"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core"
	"gx/ipfs/QmVzK524a2VWLqyvtBeiHKsUAWYgeAk4DBeZoY7vpNPNRx/go-block-format"
	"sync"
	"context"
	"math/rand"
)

//Find other peers around a subject.
//This is done by using kad-DHT.
type BitswapService struct {
	node     *core.IpfsNode
	dht      *ipdht.IpfsDHT
	mutex    sync.Mutex
}

func NewBitSwap(node *core.IpfsNode) (*BitswapService, error) {
	if dht, ok := node.Routing.(*ipdht.IpfsDHT); ok {
		return &BitswapService{
			node:     node,
			dht:      dht,
		}, nil
	} else {
		return nil, errors.New("interface conversion: node.Routing is not *ipdht.IpfsDHT")
	}
}

func (b *BitswapService) RegisterValidator(key string, validatorFunc validateFunc_t, selectFunc selectFunc_t, cache bool) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	var err error
	b.dht.Validator.(record.NamespacedValidator)[key], err = newDecentralizedValidator(validatorFunc, selectFunc)
	return err
}

func (b *BitswapService) Find(subject string, num int) <-chan pstore.PeerInfo {
	logger.Debugf("Find subject: %s", subject)
	peers := b.node.Routing.FindProvidersAsync(b.node.Context(), StringToCid(subject), num)
	logger.Debugf("Found %d around %s", len(peers), subject)
	return peers
}

func (b *BitswapService) Provide(subject string) error {
	return b.node.Routing.Provide(b.node.Context(), StringToCid(subject), true)
}

func (b *BitswapService) DecodeKey(key string) (string, error) {
	data, err := b58.Decode(key)
	return string(data), err
}

func (b *BitswapService) getKey(keyType string, rawKey string) string {
	return fmt.Sprintf("/%s/%s", keyType, b58.Encode([]byte(rawKey)))
}

func (b *BitswapService) getShardedKey(keyType string, rawKey string, slot int) string {
	return fmt.Sprintf("/%s/%s_%d", keyType, b58.Encode([]byte(rawKey)), slot)
}

func getSlot() int {
	return rand.Intn(NUM_MATCHMAKING_SLOTS)
}

func (b *BitswapService) PutValue(keyType string, rawKey string, data []byte) error {
	logger.Infof("Put value for type %s for key %s", keyType, rawKey)
	key := b.getKey(keyType, rawKey)
	return b.node.Routing.PutValue(b.node.Context(), key, data)
}

func (b *BitswapService) GetValue(ctx context.Context, keyType string, rawKey string) ([]byte, error) {
	logger.Infof("Get best value for type %s for key %s", keyType, rawKey)
	key := b.getKey(keyType, rawKey)
	return b.node.Routing.GetValue(ctx, key)
}

//Sharded value means that on this key type we can have several values. The value put in here will get shared over slots
func (b *BitswapService) PutShardedValue(keyType string, rawKey string, data []byte) error {
	logger.Infof("Put sharded value for type %s for key %s", keyType, rawKey)
	key := b.getShardedKey(keyType, rawKey, getSlot())
	return b.node.Routing.PutValue(b.node.Context(), key, data)
}

func (b *BitswapService) GetShardedValue(ctx context.Context, keyType string, rawKey string) ([][]byte, error) {
	logger.Infof("Get sharded best value for type %s for key %s", keyType, rawKey)
	var result [][]byte
	var err error
	for i := 0; i <= NUM_MATCHMAKING_SLOTS; i++ {
		key := b.getShardedKey(keyType, rawKey, i)
		var value []byte
		value, err = b.node.Routing.GetValue(ctx, key)
		if err != nil {
			logger.Debug(err)
			continue
		}
		result = append(result, value)
	}
	if len(result) == 0 {
		return nil, err
	}
	return result, nil
}

func StringToCid(value string) *cid.Cid {
	block := blocks.NewBlock([]byte(value))
	return block.Cid()
}

func StringToCid2(value string) *cid.Cid {
	return cid.NewCidV0(u.Hash([]byte(value)))
}

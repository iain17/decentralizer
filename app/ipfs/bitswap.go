package ipfs

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/golang-lru"
	"github.com/iain17/logger"
	u "gx/ipfs/QmPdKqUcHGFdeSpvjVoaTRPPstGif9GBZb5Q56RVw9o69A/go-ipfs-util"
	"gx/ipfs/QmZ383TySJVeZWzGnWui6pRcKyYZk9VkKTuW7tmKRWk5au/go-libp2p-routing"
	"gx/ipfs/QmVsp2KdPYE6M8ryzCk5KHLo3zprcY5hBDaYx6uPCFUdxA/go-libp2p-record"
	ipdht "gx/ipfs/QmTktQYCKzQjhxF6dk5xJPRuhHn3JBiKGvMLoiDy1mYmxC/go-libp2p-kad-dht"
	b58 "gx/ipfs/QmWFAMPqsEyUX7gDUsRVmMWz59FxSpJ1b2v6bJ1yYzo7jY/go-base58-fast/base58"
	pstore "gx/ipfs/QmZR2XWVVBCtbgBWnQhWk2xcQfaR3W8faQPriAiaaj7rsr/go-libp2p-peerstore"
	"gx/ipfs/QmYVNvtQkeZ6AKSwDrjQTs432QtL6umrrK41EBq3cu7iSP/go-cid"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core"
	"gx/ipfs/QmVzK524a2VWLqyvtBeiHKsUAWYgeAk4DBeZoY7vpNPNRx/go-block-format"
	"hash"
	"hash/crc32"
	"sync"
)

//Find other peers around a subject.
//This is done by using kad-DHT.
type BitswapService struct {
	node     *core.IpfsNode
	dht      *ipdht.IpfsDHT
	dhtCache *lru.Cache //Cache our result to certain DHT values.
	crcTable hash.Hash32
	mutex    sync.Mutex
}

func NewBitSwap(node *core.IpfsNode) (*BitswapService, error) {
	dhtCache, err := lru.New(4096)
	if err != nil {
		return nil, err
	}
	if dht, ok := node.Routing.(*ipdht.IpfsDHT); ok {
		return &BitswapService{
			node:     node,
			dhtCache: dhtCache,
			crcTable: crc32.NewIEEE(),
			dht:      dht,
		}, nil
	} else {
		return nil, errors.New("interface conversion: node.Routing is not *ipdht.IpfsDHT")
	}
}

func (b *BitswapService) getValidatorKey(keyType string, data []byte) uint32 {
	b.crcTable.Reset()
	b.crcTable.Write(data)
	//return fmt.Sprintf("%s/%d", keyType, b.crcTable.Sum32())
	return b.crcTable.Sum32()
}

func (b *BitswapService) RegisterValidator(keyType string, validatorFunc record.ValidatorFunc, sign bool, cache bool) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.dht.Validator[keyType] = &record.ValidChecker{
		Func: func(r *record.ValidationRecord) error {
			cacheKey := b.getValidatorKey(keyType, r.Value)
			if cacheVal, ok := b.dhtCache.Get(cacheKey); ok {
				if val, ok := cacheVal.(bool); ok {
					if val {
						return nil
					} else {
						return errors.New("cache validator return error previously")
					}
				}
			}
			result := validatorFunc(r)
			b.dhtCache.Add(cacheKey, result == nil)
			return result
		},
		Sign: sign,
	}
	if !cache {
		b.dht.Validator[keyType].Func = validatorFunc
	}
}

func (b *BitswapService) RegisterSelector(keyType string, selectorFunc record.SelectorFunc) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.dht.Selector[keyType] = selectorFunc
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

func (b *BitswapService) PutValue(keyType string, rawKey string, data []byte) error {
	logger.Infof("Put value for type %s for key %s", keyType, rawKey)
	key := b.getKey(keyType, rawKey)
	return b.node.Routing.PutValue(b.node.Context(), key, data)
	//return b.node.Routing.PutValue(b.node.Context(), key, data)
}

func (b *BitswapService) GetValues(ctx context.Context, keyType string, rawKey string, count int) ([]routing.RecvdVal, error) {
	logger.Infof("Get values for type %s for key %s", keyType, rawKey)
	key := b.getKey(keyType, rawKey)
	//return b.dht.GetValues2(b.node.Context(), key, count)//The first ifstatement, checking for cache.
	return b.node.Routing.GetValues(ctx, key, count)
}

func (b *BitswapService) GetValue(ctx context.Context, keyType string, rawKey string) ([]byte, error) {
	logger.Infof("Get best value for type %s for key %s", keyType, rawKey)
	key := b.getKey(keyType, rawKey)
	return b.node.Routing.GetValue(ctx, key)
}

func StringToCid(value string) *cid.Cid {
	block := blocks.NewBlock([]byte(value))
	return block.Cid()
}

func StringToCid2(value string) *cid.Cid {
	return cid.NewCidV0(u.Hash([]byte(value)))
}

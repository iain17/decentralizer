package ipfs

import (
	u "gx/ipfs/QmSU6eubNdhXjFBJBSksTp8kv8YRub8mGAPv8tVJHmL2EU/go-ipfs-util"
	"gx/ipfs/QmYHpXQEWuhwgRFBnrf4Ua6AZhcqXCYa7Biv65SLGgTgq5/go-ipfs/core"
	pstore "gx/ipfs/QmPgDWmTmuzvP7QE5zwo1TmjbJme9pmZHNujB2453jkCTr/go-libp2p-peerstore"
	"gx/ipfs/QmSn9Td7xgxm9EV7iEjTckpUWmWApggzPxu7eFGWkkpwin/go-block-format"
	"gx/ipfs/QmNp85zy9RLrQ5oQD4hPyS39ezrrXpcaa7R4Y9kxdWQLLQ/go-cid"
	b58 "gx/ipfs/QmT8rehPR3F6bmwL6zjUN8XpiDBFFpMP2myPdC6ApsWfJf/go-base58"
	"gx/ipfs/QmPR2JzfKd9poHx9XBhzoFeBBC31ZM3W5iUPKJZWyaoZZm/go-libp2p-routing"
	"github.com/iain17/timeout"
	"time"
	"context"
	"github.com/iain17/logger"
	ipdht "gx/ipfs/QmWRBYr99v8sjrpbyNWMuGkQekn7b9ELoLSCe8Ny7Nxain/go-libp2p-kad-dht"
	"errors"
	"fmt"
	"gx/ipfs/QmbxkgUceEcuSZ4ZdBA3x74VUDSSYjHYmmeEqkjxbtZ6Jg/go-libp2p-record"
	"strings"
	"github.com/hashicorp/golang-lru"
	"hash/crc32"
	"hash"
)

//Find other peers around a subject.
//This is done by using kad-DHT.
type BitswapService struct {
	node    *core.IpfsNode
	dht		*ipdht.IpfsDHT
	dhtCache	*lru.Cache//Cache our result to certain DHT values.
	test map[string][]byte
	crcTable hash.Hash32
}

func NewBitSwap(node *core.IpfsNode) (*BitswapService, error) {
	dhtCache, err := lru.New(4096)
	if err != nil {
		return nil, err
	}
	if dht, ok := node.Routing.(*ipdht.IpfsDHT); ok {
		return &BitswapService{
			node:    node,
			dhtCache: dhtCache,
			crcTable: crc32.NewIEEE(),
			dht: dht,
			test: make(map[string][]byte),
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

func (b *BitswapService) RegisterValidator(keyType string, validatorFunc record.ValidatorFunc, sign bool) {
	b.dht.Validator[keyType] = &record.ValidChecker{
		Func: func(key string, value[]byte) error {
			cacheKey := b.getValidatorKey(keyType, value)
			if cacheVal, ok := b.dhtCache.Get(cacheKey); ok {
				if val, ok := cacheVal.(bool); ok {
					if val {
						return nil
					} else {
						return errors.New("cache validator return error previously")
					}
				}
			}
			result := validatorFunc(key, value)
			b.dhtCache.Add(cacheKey, result == nil)
			return result
		},
		Sign: sign,
	}
}

func (b *BitswapService) RegisterSelector(keyType string, selectorFunc record.SelectorFunc) {
	b.dht.Selector[keyType] = selectorFunc
}

func (b *BitswapService) Find(subject string, num int) <-chan pstore.PeerInfo {
	logger.Debugf("Find subject: %s", subject)
	peers := b.node.Routing.FindProvidersAsync(b.node.Context(), StringToCid(subject), num)
	logger.Debugf("Found %d around %s", len(peers), subject)
	return peers
}

func (b *BitswapService) Provide(subject string) error {
	var err error
	completed := false
	timeout.Do(func(ctx context.Context) {
		err = b.node.Routing.Provide(b.node.Context(), StringToCid(subject), true)
		completed = true
		logger.Debugf("Provided subject: %s", subject)
	}, 5*time.Second)
	//if !completed {
	//	err = errors.New("could not provide '%s' in under 15 seconds. Check if you are connected to enough peers")
	//}
	return err
}

func (b *BitswapService) DecodeKey(key string) (string, error) {
	parts := strings.Split(key, "/")
	if len(parts) != 3 {
		return "", errors.New("invalid key")
	}
	return string(b58.Decode(parts[2])), nil
}

func (b *BitswapService) getKey(keyType string, rawKey string) string {
	return fmt.Sprintf("/%s/%s", keyType, b58.Encode([]byte(rawKey)))
}

func (b *BitswapService) PutValue(keyType string, rawKey string, data []byte) error {
	key := b.getKey(keyType, rawKey)
	logger.Infof("Put value: %s", key)
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

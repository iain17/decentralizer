package routinghelpers

import (
	"context"
	"strings"

	routing "gx/ipfs/QmZ383TySJVeZWzGnWui6pRcKyYZk9VkKTuW7tmKRWk5au/go-libp2p-routing"
	ropts "gx/ipfs/QmZ383TySJVeZWzGnWui6pRcKyYZk9VkKTuW7tmKRWk5au/go-libp2p-routing/options"

	peer "gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	ci "gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"
)

// LimitedValueStore limits the internal value store to the given namespaces.
type LimitedValueStore struct {
	routing.ValueStore
	Namespaces []string
}

// GetPublicKey returns the public key for the given peer.
func (lvs *LimitedValueStore) GetPublicKey(ctx context.Context, p peer.ID) (ci.PubKey, error) {
	for _, ns := range lvs.Namespaces {
		if ns == "pk" {
			return routing.GetPublicKey(lvs.ValueStore, ctx, p)
		}
	}
	return nil, routing.ErrNotFound
}

// PutValue returns ErrNotSupported
func (lvs *LimitedValueStore) PutValue(ctx context.Context, key string, value []byte, opts ...ropts.Option) error {
	if !lvs.KeySupported(key) {
		return routing.ErrNotSupported
	}
	return lvs.ValueStore.PutValue(ctx, key, value, opts...)
}

// KeySupported returns true if the passed key is supported by this value store.
func (lvs *LimitedValueStore) KeySupported(key string) bool {
	if len(key) < 3 {
		return false
	}
	if key[0] != '/' {
		return false
	}
	key = key[1:]
	for _, ns := range lvs.Namespaces {
		if len(ns) < len(key) && strings.HasPrefix(key, ns) && key[len(ns)] == '/' {
			return true
		}
	}
	return false
}

// GetValue returns ErrNotSupported
func (lvs *LimitedValueStore) GetValue(ctx context.Context, key string, opts ...ropts.Option) ([]byte, error) {
	if !lvs.KeySupported(key) {
		return nil, routing.ErrNotFound
	}
	return lvs.ValueStore.GetValue(ctx, key, opts...)
}

func (lvs *LimitedValueStore) Bootstrap(ctx context.Context) error {
	if bs, ok := lvs.ValueStore.(Bootstrap); ok {
		return bs.Bootstrap(ctx)
	}
	return nil
}

var _ routing.PubKeyFetcher = (*LimitedValueStore)(nil)
var _ routing.ValueStore = (*LimitedValueStore)(nil)
var _ Bootstrap = (*LimitedValueStore)(nil)

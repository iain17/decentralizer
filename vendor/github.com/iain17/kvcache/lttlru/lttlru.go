// Package lttlru LRU with TTL, implemented without goroutine.
package lttlru

import (
	"time"
	"sync"

	"github.com/pkg/errors"
	hlru "github.com/hashicorp/golang-lru"
)

// LruWithTTL lru with ttl
type LruWithTTL struct {
	*hlru.Cache
	expiresAt map[interface{}]*time.Time
	expireMutex sync.Mutex
}

// NewTTL creates an LRU of the given size
func NewTTL(size int) (*LruWithTTL, error) {
	return NewTTLWithEvict(size, nil)
}

// NewTTLWithEvict creates an LRU of the given size and a evict callback function
func NewTTLWithEvict(size int, onEvicted func(key interface{}, value interface{})) (*LruWithTTL, error) {
	if size <= 0 {
		return nil, errors.New("Must provide a positive size")
	}
	c, err := hlru.NewWithEvict(size, onEvicted)
	if err != nil {
		return nil, err
	}
	return &LruWithTTL{c, make(map[interface{}]*time.Time), sync.Mutex{}}, nil
}

func (lru *LruWithTTL) removeExpired(key interface{}) {
	lru.expireMutex.Lock()
	defer lru.expireMutex.Unlock()
	delete(lru.expiresAt, key)
	lru.Remove(key)
}

// AddWithTTL add an key:val with TTL
func (lru *LruWithTTL) AddWithTTL(key, value interface{}, ttl time.Duration) bool {
	lru.expireMutex.Lock()
	defer lru.expireMutex.Unlock()
	expire := time.Now().Add(ttl)
	lru.expiresAt[key] = &expire
	return lru.Cache.Add(key, value)
}

func (lru *LruWithTTL) GetWithTTL(key interface{}) (interface{}, bool) {
	if lru.expiresAt[key] != nil {
		if lru.expiresAt[key].After(time.Now()) {
			// ttl not expired
			return lru.Get(key)
		}

		// ttl expired
		lru.removeExpired(key)
	}

	// not added with ttl
	return nil, false
}


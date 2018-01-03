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
	expireMutex sync.RWMutex
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
	return &LruWithTTL{c, make(map[interface{}]*time.Time), sync.RWMutex{}}, nil
}

func (lru *LruWithTTL) removeExpired(key interface{}) {
	lru.expireMutex.Lock()
	delete(lru.expiresAt, key)
	lru.Remove(key)
	lru.expireMutex.Unlock()
}

// AddWithTTL add an key:val with TTL
func (lru *LruWithTTL) AddWithTTL(key, value interface{}, ttl time.Duration) bool {
	lru.expireMutex.Lock()
	expire := time.Now().Add(ttl)
	lru.expiresAt[key] = &expire
	lru.expireMutex.Unlock()

	return lru.Cache.Add(key, value)
}

func (lru *LruWithTTL) Contains(key interface{}) bool {
	lru.expireMutex.RLock()
	if lru.expiresAt[key] != nil {
		if !lru.expiresAt[key].After(time.Now()) {
			lru.removeExpired(key)
			lru.expireMutex.RUnlock()
			return false
		}
	}
	lru.expireMutex.RUnlock()
	return lru.Cache.Contains(key)
}

func (lru *LruWithTTL) Peek(key interface{}) (interface{}, bool) {
	lru.expireMutex.RLock()
	if lru.expiresAt[key] != nil {
		if !lru.expiresAt[key].After(time.Now()) {
			lru.removeExpired(key)
			lru.expireMutex.RUnlock()
			return nil, false
		}
	}
	lru.expireMutex.RUnlock()
	return lru.Cache.Peek(key)
}

func (lru *LruWithTTL) Get(key interface{}) (interface{}, bool) {
	lru.expireMutex.RLock()
	if lru.expiresAt[key] != nil {
		if !lru.expiresAt[key].After(time.Now()) {
			lru.removeExpired(key)
			lru.expireMutex.RUnlock()
			return nil, false
		}
	}
	lru.expireMutex.RUnlock()
	return lru.Cache.Get(key)
}


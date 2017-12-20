// Package ttlru LRU with TTL, implemented with goroutine.
package lru

import (
	"errors"
	"sync"
	"time"

	hlru "github.com/hashicorp/golang-lru"
)

// LruWithTTL lru with ttl
type LruWithTTL struct {
	*hlru.Cache
	schedule      map[interface{}]*time.Timer
	scheduleMutex sync.Mutex
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
	return &LruWithTTL{c, make(map[interface{}]*time.Timer), sync.Mutex{}}, nil
}

func (lru *LruWithTTL) clearSchedule(key interface{}) {
	lru.scheduleMutex.Lock()
	defer lru.scheduleMutex.Unlock()
	delete(lru.schedule, key)
}

// AddWithTTL add an key:val with TTL
func (lru *LruWithTTL) AddWithTTL(key, value interface{}, ttl time.Duration) bool {
	lru.scheduleMutex.Lock()
	defer lru.scheduleMutex.Unlock()
	if lru.schedule[key] != nil {
		// already scheduled, nothing to do
		lru.schedule[key].Reset(ttl)
	} else {
		lru.schedule[key] = time.NewTimer(ttl)
		// Schedule cleanup
		go func() {
			<-lru.schedule[key].C
			lru.Cache.Remove(key)
			lru.clearSchedule(key)
		}()
	}

	return lru.Cache.Add(key, value)
}

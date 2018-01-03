// Package ttlru LRU with TTL, implemented with goroutine.
package lru

import (
	"errors"
	"sync"
	"time"

	hlru "github.com/hashicorp/golang-lru"
	"context"
)

// LruWithTTL lru with ttl
type LruWithTTL struct {
	*hlru.Cache
	schedule      map[interface{}]time.Time
	scheduleMutex sync.RWMutex
}

// NewTTL creates an LRU of the given size
func NewTTL(contex context.Context, size int) (*LruWithTTL, error) {
	return NewTTLWithEvict(contex, size, nil)
}

// NewTTLWithEvict creates an LRU of the given size and a evict callback function
func NewTTLWithEvict(context context.Context, size int, onEvicted func(key interface{}, value interface{})) (*LruWithTTL, error) {
	if size <= 0 {
		return nil, errors.New("Must provide a positive size")
	}
	c, err := hlru.NewWithEvict(size, onEvicted)
	if err != nil {
		return nil, err
	}

	i:= &LruWithTTL{c, make(map[interface{}]time.Time), sync.RWMutex{}}
	go i.runSchedule(context)
	return i, nil
}

func (lru *LruWithTTL) runSchedule(context context.Context) {
	timer := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		case <-context.Done():
			return
		case <- timer.C:
			lru.scheduleMutex.Lock()
			now := time.Now()
			for key, ttl := range lru.schedule {
				if !now.Before(ttl) {
					lru.Cache.Remove(key)
					delete(lru.schedule, key)
				}
			}
			lru.scheduleMutex.Unlock()
		}
	}
}

// AddWithTTL add an key:val with TTL
func (lru *LruWithTTL) AddWithTTL(key, value interface{}, ttl time.Duration) bool {
	lru.scheduleMutex.Lock()
	lru.schedule[key] = time.Now().Add(ttl)
	lru.scheduleMutex.Unlock()

	return lru.Cache.Add(key, value)
}

package lru

import (
	"testing"
	"time"
	"context"
)

// test that Add returns true/false if an eviction occurred
func TestLRUTTLAddNoTTL(t *testing.T) {
	evictCounter := 0
	onEvicted := func(k interface{}, v interface{}) {
		evictCounter += 1
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	l, err := NewTTLWithEvict(ctx, 1, onEvicted)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if l.Add(1, 1) || evictCounter != 0 {
		t.Errorf("should not have an eviction")
	}
	if !l.Add(2, 2) || evictCounter != 1 {
		t.Errorf("should have an eviction")
	}
}

// test that Add returns true/false if an eviction occurred
func TestLRUTTLAddWithTTL(t *testing.T) {
	evictCounter := 0
	onEvicted := func(k interface{}, v interface{}) {
		evictCounter++
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	l, err := NewTTLWithEvict(ctx, 2, onEvicted)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if l.AddWithTTL(1, 1, time.Millisecond*5) {
		t.Errorf("should not have an eviction")
	}
	if l.AddWithTTL(2, 2, time.Millisecond*10) {
		t.Errorf("should have an eviction")
	}

	// Wait for TTLs to expire
	time.Sleep(200 * time.Millisecond)

	if evictCounter != 2 {
		t.Errorf("should have been 2 evictions")
	}
}

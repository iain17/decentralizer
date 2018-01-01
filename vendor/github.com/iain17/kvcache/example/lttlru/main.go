package main

import (
	"fmt"
	"time"

	"github.com/Akagi201/kvcache/lttlru"
)

func main() {
	evictCounter := 0
	onEvicted := func(k interface{}, v interface{}) {
		fmt.Printf("expired, k: %v, v: %v\n", k, v)
		evictCounter += 1
		if v.(int) != evictCounter {
			fmt.Errorf("eviction happened out of order. Got %v, expected %v", v.(int), evictCounter)
		}
	}

	l, err := lttlru.NewTTLWithEvict(2, onEvicted)
	if err != nil {
		fmt.Errorf("err: %v\n", err)
		return
	}

	// ttl expired test
	if l.AddWithTTL(1, 1, time.Second*2) {
		fmt.Errorf("should not have an eviction\n")
	}

	time.Sleep(1 * time.Second)
	v, b := l.GetWithTTL(1)
	fmt.Printf("v: %v, b: %v\n", v, b)

	time.Sleep(2 * time.Second)
	v, b = l.GetWithTTL(1)
	fmt.Printf("v: %v, b: %v\n", v, b)

	// ttl expire update test
	if l.AddWithTTL(1, 1, time.Second*3) {
		fmt.Errorf("should not have an eviction\n")
	}

	time.Sleep(2 * time.Second)
	if l.AddWithTTL(1, 1, time.Second*3) {
		fmt.Errorf("should not have an eviction\n")
	}
	v, b = l.GetWithTTL(1)
	fmt.Printf("v: %v, b: %v\n", v, b)

	time.Sleep(2 * time.Second)
	if l.AddWithTTL(1, 1, time.Second*3) {
		fmt.Errorf("should not have an eviction\n")
	}
	v, b = l.GetWithTTL(1)
	fmt.Printf("v: %v, b: %v\n", v, b)

	fmt.Printf("evictCounter: %v\n", evictCounter)
}

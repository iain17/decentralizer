package main

import (
	"fmt"
	"time"

	"github.com/Akagi201/kvcache/ttlru"
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

	l, err := lru.NewTTLWithEvict(2, onEvicted)
	if err != nil {
		fmt.Errorf("err: %v\n", err)
		return
	}

	if l.AddWithTTL(1, 1, time.Second*5) {
		fmt.Errorf("should not have an eviction\n")
	}

	if l.AddWithTTL(2, 2, time.Second*10) {
		fmt.Errorf("should have an eviction\n")
	}

	time.Sleep(2 * time.Second)
	v, b := l.Get(1)
	fmt.Printf("v: %v, b: %v\n", v, b)

	time.Sleep(2 * time.Second)
	v, b = l.Get(1)
	if l.AddWithTTL(1, 1, time.Second*5) {
		fmt.Errorf("should not have an eviction\n")
	}
	fmt.Printf("v: %v, b: %v\n", v, b)

	time.Sleep(2 * time.Second)
	v, b = l.Get(1)
	if l.AddWithTTL(1, 1, time.Second*5) {
		fmt.Errorf("should not have an eviction\n")
	}
	fmt.Printf("v: %v, b: %v\n", v, b)

	time.Sleep(2 * time.Second)
	v, b = l.Get(1)
	fmt.Printf("v: %v, b: %v\n", v, b)

	time.Sleep(5 * time.Second)
	v, b = l.Get(1)
	fmt.Printf("v: %v, b: %v\n", v, b)

	fmt.Printf("evictCounter: %v\n", evictCounter)
}

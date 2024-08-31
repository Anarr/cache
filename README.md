# Go Cache Library

This Go library provides implementations of three different types of caches:

1. **Least Recently Used (LRU) Cache**
2. **Least Frequently Used (LFU) Cache**
3. **Time-Based Expiration Cache**

These caches can be used to efficiently manage data with various eviction strategies based on your application's needs.

## Installation

To install the Go Cache Library, use the `go get` command:

```bash
go get github.com/Anarr/cache

package main

import (
	"fmt"
	"gocache/cache"
)

func main() {
	// create an LFU cache with a capacity of 2
	lfuCache := cache.NewLFUCache(2)

	// add items to the cache
	lfuCache.Put("a", 10)
	lfuCache.Put("b", 20)

	// Retrieve an item
	val, _ := lfuCache.Get("a")
	fmt.Println("Value for key 'a':", val) // Output: 10

	// access the key "a" again to increase its frequency
	lfuCache.Get("a")

	// Add a new item, which will cause the least frequently used item ('b') to be evicted
	lfuCache.Put("c", 30)

	// check if the least frequently used item ('b') was evicted
	_, found := lfuCache.Get("b")
	fmt.Println("Key 'b' exists:", found) // Output: false

	// check if 'a' and 'c' are still in the cache
	val, _ = lfuCache.Get("a")
	fmt.Println("Value for key '


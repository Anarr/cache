package cache

import (
	"container/list"
	"sync"
)

// LRUCache implements a Least Recently Used cache.
type LRUCache struct {
	capacity int
	items    map[string]*list.Element
	order    *list.List
	mu       sync.RWMutex
}

// Entry represents a cache entry for LRU.
type Entry struct {
	key   string
	value any
}

// NewLRUCache initializes a new LRU cache with a specified capacity.
func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		items:    make(map[string]*list.Element),
		order:    list.New(),
	}
}

// Get retrieves an item from the LRU cache.
func (c *LRUCache) Get(key string) (any, bool) {
	c.mu.RLock()
	element, found := c.items[key]
	c.mu.RUnlock()

	if !found {
		return nil, false
	}

	c.mu.Lock()
	c.order.MoveToFront(element)
	c.mu.Unlock()

	return element.Value.(*Entry).value, true
}

// Put inserts an item into the LRU cache.
func (c *LRUCache) Put(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if element, found := c.items[key]; found {
		c.order.MoveToFront(element)
		element.Value.(*Entry).value = value
		return
	}

	if len(c.items) >= c.capacity {
		c.evict()
	}

	entry := &Entry{key, value}
	element := c.order.PushFront(entry)
	c.items[key] = element
}

// Remove deletes an item from the LRU cache.
func (c *LRUCache) Remove(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if element, found := c.items[key]; found {
		c.order.Remove(element)
		delete(c.items, key)
	}
}

// evict removes the least recently used item from the cache.
func (c *LRUCache) evict() {
	element := c.order.Back()
	if element != nil {
		c.order.Remove(element)
		entry := element.Value.(*Entry)
		delete(c.items, entry.key)
	}
}

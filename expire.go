package cache

import (
	"sync"
	"time"
)

// TimeBasedCache implements a cache with time-based expiration
type TimeBasedCache struct {
	capacity   int
	items      map[string]*Item
	expiration time.Duration
	mu         sync.RWMutex
}

// NewTimeBasedCache initializes a new time-based cache
func NewTimeBasedCache(capacity int, expiration time.Duration) *TimeBasedCache {
	return &TimeBasedCache{
		capacity:   capacity,
		items:      make(map[string]*Item),
		expiration: expiration,
	}
}

// Get retrieves an item from the cache and checks for expiration
func (c *TimeBasedCache) Get(key string) (any, bool) {
	c.mu.RLock()
	item, found := c.items[key]
	c.mu.RUnlock()

	if !found || time.Since(item.timestamp) > c.expiration {
		c.Remove(key)
		return nil, false
	}

	return item.value, true
}

// Put inserts an item into the cache
func (c *TimeBasedCache) Put(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.items) >= c.capacity {
		c.evict()
	}

	item := &Item{key: key, value: value, timestamp: time.Now()}
	c.items[key] = item
}

// Remove deletes an item from the cache
func (c *TimeBasedCache) Remove(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// evict removes an expired item from the cache.
func (c *TimeBasedCache) evict() {
	var oldestKey string
	var oldestTime time.Time

	for key, item := range c.items {
		if oldestKey == "" || item.timestamp.Before(oldestTime) {
			oldestKey = key
			oldestTime = item.timestamp
		}
	}

	if oldestKey != "" {
		delete(c.items, oldestKey)
	}
}

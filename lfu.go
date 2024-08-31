package cache

import (
	"sync"
	"time"
)

// Item represents a cache item, used in LFU and Time-Based caches
type Item struct {
	key       string
	value     any
	frequency int       // Used in LFU cache
	timestamp time.Time // Used in Time-Based cache
	index     int       // Index in the heap
}

// LFUCache implements a Least Frequently Used cache
type LFUCache struct {
	capacity int
	items    map[string]*Item
	freqHeap []*Item
	mu       sync.RWMutex
}

// NewLFUCache initializes a new LFU cache with a specified capacity
func NewLFUCache(capacity int) *LFUCache {
	return &LFUCache{
		capacity: capacity,
		items:    make(map[string]*Item),
		freqHeap: make([]*Item, 0, capacity),
	}
}

// Get retrieves an item from the LFU cache
func (c *LFUCache) Get(key string) (any, bool) {
	c.mu.RLock()
	item, found := c.items[key]
	c.mu.RUnlock()

	if !found {
		return nil, false
	}

	c.mu.Lock()
	item.frequency++
	c.fixHeap(item.index)
	c.mu.Unlock()

	return item.value, true
}

// Put inserts an item into the LFU cache
func (c *LFUCache) Put(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, found := c.items[key]; found {
		item.value = value
		item.frequency++
		c.fixHeap(item.index)
		return
	}

	if len(c.items) >= c.capacity {
		c.evict()
	}

	item := &Item{key: key, value: value, frequency: 1}
	c.items[key] = item
	c.pushHeap(item)
}

// Remove deletes an item from the LFU cache
func (c *LFUCache) Remove(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, found := c.items[key]; found {
		c.removeHeap(item.index)
		delete(c.items, key)
	}
}

// pushHeap adds a new item to the heap
func (c *LFUCache) pushHeap(item *Item) {
	item.index = len(c.freqHeap)
	c.freqHeap = append(c.freqHeap, item)
	c.upHeap(len(c.freqHeap) - 1)
}

// removeHeap removes an item from the heap by index
func (c *LFUCache) removeHeap(index int) {
	if index >= len(c.freqHeap) {
		return
	}
	last := len(c.freqHeap) - 1
	c.swapHeap(index, last)
	c.freqHeap = c.freqHeap[:last]
	c.downHeap(index)
}

// fixHeap reorders the heap after an items frequency is updated
func (c *LFUCache) fixHeap(index int) {
	c.downHeap(index)
	c.upHeap(index)
}

// evict removes the least frequently used item from the cache.
func (c *LFUCache) evict() {
	if len(c.freqHeap) > 0 {
		item := c.freqHeap[0]
		c.removeHeap(0)
		delete(c.items, item.key)
	}
}

// upHeap restores the heap property by moving an item up
func (c *LFUCache) upHeap(index int) {
	for index > 0 {
		parent := (index - 1) / 2
		if c.freqHeap[index].frequency >= c.freqHeap[parent].frequency {
			break
		}
		c.swapHeap(index, parent)
		index = parent
	}
}

// downHeap restores the heap property by moving an item down
func (c *LFUCache) downHeap(index int) {
	for {
		left := 2*index + 1
		right := 2*index + 2
		smallest := index

		if left < len(c.freqHeap) && c.freqHeap[left].frequency < c.freqHeap[smallest].frequency {
			smallest = left
		}
		if right < len(c.freqHeap) && c.freqHeap[right].frequency < c.freqHeap[smallest].frequency {
			smallest = right
		}
		if smallest == index {
			break
		}
		c.swapHeap(index, smallest)
		index = smallest
	}
}

// swapHeap swaps two items in the heap
func (c *LFUCache) swapHeap(i, j int) {
	c.freqHeap[i], c.freqHeap[j] = c.freqHeap[j], c.freqHeap[i]
	c.freqHeap[i].index = i
	c.freqHeap[j].index = j
}

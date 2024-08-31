package cache_test

import (
	"github.com/Anarr/cache"
	"testing"
)

func TestLRUCache(t *testing.T) {
	cache := cache.NewLRUCache(2)

	// Test Put and Get functionality.
	cache.Put("a", 1)
	cache.Put("b", 2)

	if val, ok := cache.Get("a"); !ok || val != 1 {
		t.Errorf("Expected value 1 for key 'a', got %v", val)
	}
	if val, ok := cache.Get("b"); !ok || val != 2 {
		t.Errorf("Expected value 2 for key 'b', got %v", val)
	}

	// Adding a new item should evict the least recently used item ('a')
	cache.Put("c", 3)
	if _, ok := cache.Get("a"); ok {
		t.Errorf("Expected key 'a' to be evicted")
	}

	// 'b' should still be present and retrievable.
	if val, ok := cache.Get("b"); !ok || val != 2 {
		t.Errorf("Expected value 2 for key 'b', got %v", val)
	}
}

func TestLRUCacheEviction(t *testing.T) {
	cache := cache.NewLRUCache(2)
	cache.Put("a", 1)
	cache.Put("b", 2)

	cache.Get("a")

	// Adding a new item should evict 'b' (least recently used).
	cache.Put("c", 3)
	if _, ok := cache.Get("b"); ok {
		t.Errorf("Expected key 'b' to be evicted")
	}
}

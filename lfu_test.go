package cache_test

import (
	"github.com/Anarr/cache"
	"testing"
)

func TestLFUCache(t *testing.T) {
	cache := cache.NewLFUCache(2)

	cache.Put("x", 10)
	cache.Put("y", 20)

	if val, ok := cache.Get("x"); !ok || val != 10 {
		t.Errorf("Expected value 10 for key 'x', got %v", val)
	}
	if val, ok := cache.Get("y"); !ok || val != 20 {
		t.Errorf("Expected value 20 for key 'y', got %v", val)
	}

	cache.Get("x")

	// Adding a new item should evict the least frequently used item ('y')
	cache.Put("z", 30)
	if _, ok := cache.Get("y"); ok {
		t.Errorf("Expected key 'y' to be evicted")
	}

	// 'x' should still be present and retrievable.
	if val, ok := cache.Get("x"); !ok || val != 10 {
		t.Errorf("Expected value 10 for key 'x', got %v", val)
	}
}

func TestLFUCacheEviction(t *testing.T) {
	cache := cache.NewLFUCache(2)
	cache.Put("a", 1)
	cache.Put("b", 2)

	// Access a to make it more frequently used
	cache.Get("a")

	// Adding a new item should evict 'b' (less frequently used)
	cache.Put("c", 3)
	if _, ok := cache.Get("b"); ok {
		t.Errorf("Expected key 'b' to be evicted")
	}

	// 'a' should still be present and retrievable
	if val, ok := cache.Get("a"); !ok || val != 1 {
		t.Errorf("Expected value 1 for key 'a', got %v", val)
	}
}

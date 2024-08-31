package cache_test

import (
	"github.com/Anarr/cache"
	"testing"
	"time"
)

func TestTimeBasedCache(t *testing.T) {
	cache := cache.NewTimeBasedCache(2, 2*time.Second)
	cache.Put("k", 100)

	// wait for 1 second and verify that value is still available
	time.Sleep(1 * time.Second)
	if val, ok := cache.Get("k"); !ok || val != 100 {
		t.Errorf("Expected value 100 for key 'k', got %v", val)
	}

	// Wait for another 2 seconds to exceed the expiration time.
	time.Sleep(2 * time.Second)

	// Verify that the value has expired and is no longer retrievable.
	if _, ok := cache.Get("k"); ok {
		t.Errorf("Expected key 'k' to be expired and evicted")
	}
}

func TestTimeBasedCacheEviction(t *testing.T) {
	cache := cache.NewTimeBasedCache(2, 1*time.Second)

	cache.Put("a", "valueA")
	cache.Put("b", "valueB")

	if val, ok := cache.Get("a"); !ok || val != "valueA" {
		t.Errorf("Expected value 'valueA' for key 'a', got %v", val)
	}
	if val, ok := cache.Get("b"); !ok || val != "valueB" {
		t.Errorf("Expected value 'valueB' for key 'b', got %v", val)
	}

	// wait for 1 second to expire the items
	time.Sleep(1 * time.Second)

	// both items should be expired
	if _, ok := cache.Get("a"); ok {
		t.Errorf("Expected key 'a' to be expired and evicted")
	}
	if _, ok := cache.Get("b"); ok {
		t.Errorf("Expected key 'b' to be expired and evicted")
	}
}

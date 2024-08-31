package cache

// Cache is the common interface for all cache implementations
type Cache interface {
	Get(key string) (any, bool)
	Put(key string, value any)
	Remove(key string)
}

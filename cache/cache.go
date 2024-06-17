package cache

import (
	"github.com/auvitly/go-tools/cache/internal"
	"sync"
)

type Cache[K comparable, V any] struct {
	storage map[K]internal.Item[V]
	cfg     config
	mu      sync.Mutex
}

// New - creating a cache instance with options.
func New[K comparable, V any](options ...Option) *Cache[K, V] {
	var cfg config

	for _, option := range options {
		option(&cfg)
	}

	return &Cache[K, V]{
		storage: make(map[K]internal.Item[V]),
		cfg:     cfg,
	}
}

// Get - getting value by key.
func (c *Cache[K, V]) Get(key K) (value V, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if result, ok := c.storage[key]; ok {
		for _, fn := range result.Expirations {
			if fn.IsExpired() {
				delete(c.storage, key)

				return value, false
			}
		}

		return result.Value, ok
	}

	return value, false
}

// Set - setting value by key.
func (c *Cache[K, V]) Set(key K, value V, options ...Option) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.storage == nil {
		c.storage = make(map[K]internal.Item[V])
	}

	var (
		item internal.Item[V]
		cfg  = c.cfg
	)

	item.Value = value

	for _, option := range options {
		option(&cfg)
	}

	item.Expirations = cfg.getExpirations()

	c.storage[key] = item
}

// Delete - delete value by key.
func (c *Cache[K, V]) Delete(keys ...K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, key := range keys {
		delete(c.storage, key)
	}
}

// Clean - clear cache.
func (c *Cache[K, V]) Clean() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.storage = make(map[K]internal.Item[V])
}

// Snapshot - take a snapshot of data from the cache.
func (c *Cache[K, V]) Snapshot() map[K]V {
	c.mu.Lock()
	defer c.mu.Unlock()

	var result = make(map[K]V)

	for key, item := range c.storage {
		result[key] = item.Value
	}

	return result
}

package inmemory

import (
	"context"
	"sync"
	"time"

	"github.com/auvitly/go-tools/standard/collection/cache"
)

// Cache - contains in memory storage that allows concurrent writing and reading.
type Cache[K comparable, V any] struct {
	storage map[K]*cache.Item[V]
	config  Config
	mu      sync.RWMutex
}

// New - creating a cache instance with options.
func New[K comparable, V any](config Config) *Cache[K, V] {
	return &Cache[K, V]{
		storage: make(map[K]*cache.Item[V]),
		config:  config,
	}
}

// Lookup - getting value by key.
func (c *Cache[K, V]) Lookup(_ context.Context, key K) (cache.Item[V], bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.storage[key]

	return *item, ok
}

// Get - getting value by key.
func (c *Cache[K, V]) Get(ctx context.Context, key K) cache.Item[V] {
	item, _ := c.Lookup(ctx, key)

	return item
}

// Set - setting value by key.
func (c *Cache[K, V]) Set(_ context.Context, key K, value cache.Item[V]) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.storage == nil {
		c.storage = make(map[K]*cache.Item[V])
	}

	c.storage[key] = &cache.Item[V]{
		Value: value.Value,
		Deadline: func() *time.Time {
			switch {
			case value.Deadline != nil:
				return value.Deadline
			case c.config.RecordLifeTime != 0:
				var deadline = time.Now().Add(c.config.RecordLifeTime)

				return &deadline
			default:
				return nil
			}
		}(),
	}
}

// Delete - delete value by key.
func (c *Cache[K, V]) Delete(_ context.Context, keys ...K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, key := range keys {
		delete(c.storage, key)
	}
}

// GC - clear cache.
func (c *Cache[K, V]) GC(_ context.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, value := range c.storage {
		if value.Deadline == nil {
			continue
		}

		if time.Until(*value.Deadline) < 0 {
			delete(c.storage, key)
		}
	}

	if c.config.RecordLimit == 0 {
		return
	}

	if len(c.storage) < c.config.RecordLimit {
		return
	}

	for key := range c.storage {
		delete(c.storage, key)

		if len(c.storage) == c.config.RecordLimit {
			break
		}
	}
}

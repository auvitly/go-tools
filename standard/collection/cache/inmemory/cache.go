package inmemory

import (
	"context"
	"sync"
	"time"

	"github.com/auvitly/go-tools/standard/collection/cache"
	"github.com/auvitly/go-tools/stderrs"
)

// Cache - contains in memory storage that allows concurrent writing and reading.
type Cache[K comparable, V any] struct {
	storage map[K]item[V]
	config  Config
	mu      sync.RWMutex
}

type item[V any] struct {
	Value    V
	Deadline *time.Time
}

// New - creating a cache instance with options.
func New[K comparable, V any](config Config) *Cache[K, V] {
	return &Cache[K, V]{
		storage: make(map[K]item[V]),
		config:  config,
	}
}

// Get - getting value by key.
func (c *Cache[K, V]) Get(_ context.Context, key K) (V, *stderrs.Error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.storage[key]
	if !ok {
		return item.Value, stderrs.NotFound.
			SetMessage("not found value in cache with key %v", key)
	}

	return item.Value, nil
}

// Set - setting value by key.
func (c *Cache[K, V]) Set(_ context.Context, key K, value V, options cache.Options) *stderrs.Error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.storage == nil {
		c.storage = make(map[K]item[V])
	}

	c.storage[key] = item[V]{
		Value: value,
		Deadline: func() *time.Time {
			switch {
			case options.TTL != nil:
				var deadline = time.Now().Add(*options.TTL)

				return &deadline
			case c.config.DefaultTTL != 0:
				var deadline = time.Now().Add(c.config.DefaultTTL)

				return &deadline
			default:
				return nil
			}
		}(),
	}

	return nil
}

// Delete - delete value by key.
func (c *Cache[K, V]) Delete(_ context.Context, keys ...K) *stderrs.Error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, key := range keys {
		delete(c.storage, key)
	}

	return nil
}

// GC - clear cache.
func (c *Cache[K, V]) GC() *stderrs.Error {
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
		return nil
	}

	if len(c.storage) < c.config.RecordLimit {
		return nil
	}

	for key := range c.storage {
		delete(c.storage, key)

		if len(c.storage) == c.config.RecordLimit {
			break
		}
	}

	return nil
}

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
	storage map[K]cache.Item[V]
	config  Config
	mu      sync.RWMutex
}

// New - creating a cache instance with options.
func New[K comparable, V any](config Config) *Cache[K, V] {
	return &Cache[K, V]{
		storage: make(map[K]cache.Item[V]),
		config:  config,
	}
}

// Get - getting value by key.
func (c *Cache[K, V]) Get(_ context.Context, key K) (cache.Item[V], *stderrs.Error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.storage[key]
	if !ok {
		return cache.Item[V]{}, stderrs.NotFound.
			SetMessage("not found value in cache with key %v", key)
	}

	return item, nil
}

// Set - setting value by key.
func (c *Cache[K, V]) Set(_ context.Context, key K, value cache.Item[V]) *stderrs.Error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.storage == nil {
		c.storage = make(map[K]cache.Item[V])
	}

	c.storage[key] = cache.Item[V]{
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

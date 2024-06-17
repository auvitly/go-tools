package cache

import (
	"github.com/auvitly/go-tools/cache/internal"
	"sync"
)

type Cache[K comparable, V any] struct {
	storage map[K]internal.Item[V]
	config  Config
	mu      sync.Mutex
}

func New[K comparable, V any](options ...Option) *Cache[K, V] {
	var config Config

	for _, option := range options {
		option(&config)
	}

	return &Cache[K, V]{
		storage: make(map[K]internal.Item[V]),
		config:  config,
	}
}

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

func (c *Cache[K, V]) Set(key K, value V, options ...Option) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.storage == nil {
		c.storage = make(map[K]internal.Item[V])
	}

	var (
		item   internal.Item[V]
		config = c.config
	)

	item.Value = value

	for _, option := range options {
		option(&config)
	}

	item.Expirations = config.getExpirations()

	c.storage[key] = item
}

func (c *Cache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.storage, key)
}

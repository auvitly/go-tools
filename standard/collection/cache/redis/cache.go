package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/auvitly/go-tools/standard/collection/cache"
	"github.com/auvitly/go-tools/stderrs"
	redis_cache "github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

// Cache - contains in memory storage that allows concurrent writing and reading.
type Cache[K comparable, V any] struct {
	storage *redis_cache.Cache
	config  Config
}

func New[K comparable, V any](config Config) (*Cache[K, V], *stderrs.Error) {
	return &Cache[K, V]{
		config: config,
		storage: redis_cache.New(&redis_cache.Options{
			Redis:      redis.NewRing(&config.Redis),
			LocalCache: config.LocalCache,
		}),
	}, nil
}

func (c *Cache[K, V]) Set(key K, value cache.Item[V]) *stderrs.Error {
	return c.SetWithContext(context.Background(), key, value)
}

func (c *Cache[K, V]) SetWithContext(ctx context.Context, key K, value cache.Item[V]) *stderrs.Error {
	err := c.storage.Set(&redis_cache.Item{
		Ctx:            context.Background(),
		Key:            fmt.Sprintf("%v", key),
		Value:          value,
		SkipLocalCache: c.config.LocalCache == nil,
		TTL: func() time.Duration {
			if value.Deadline != nil {
				return time.Until(*value.Deadline)
			}

			return 0
		}(),
	})
	if err != nil {
		return stderrs.Internal.SetMessage(err.Error())
	}

	return nil
}

func (c *Cache[K, V]) Get(key K) (cache.Item[V], *stderrs.Error) {
	return c.GetWithContext(context.Background(), key)
}

func (c *Cache[K, V]) GetWithContext(ctx context.Context, key K) (cache.Item[V], *stderrs.Error) {
	var value cache.Item[V]

	err := c.storage.Get(ctx, fmt.Sprintf("%v", key), &value.Value)
	if err != nil {
		return value, stderrs.Internal.SetMessage(err.Error())
	}

	return value, nil
}

func (c *Cache[K, V]) Delete(keys ...K) *stderrs.Error {
	return c.DeleteWithContext(context.Background(), keys...)
}

func (c *Cache[K, V]) DeleteWithContext(ctx context.Context, keys ...K) *stderrs.Error {
	for _, key := range keys {
		err := c.storage.Delete(ctx, fmt.Sprintf("%v", key))
		if err != nil {
			return stderrs.Internal.SetMessage(err.Error())
		}
	}

	return nil
}

func (c *Cache[K, V]) GC() *stderrs.Error {
	return nil
}

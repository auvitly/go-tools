package inredis

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
			Marshal:    config.MarshalFunc,
			Unmarshal:  config.UnmarshalFunc,
		}),
	}, nil
}

func (c *Cache[K, V]) Set(ctx context.Context, key K, value V, options cache.Options) *stderrs.Error {
	err := c.storage.Set(&redis_cache.Item{
		Ctx:            ctx,
		Key:            fmt.Sprintf("%v", key),
		Value:          value,
		SkipLocalCache: c.config.LocalCache == nil,
		TTL: func() time.Duration {
			switch {
			case options.TTL != nil:
				return *options.TTL
			default:
				return c.config.DefaultTTL
			}
		}(),
	})
	if err != nil {
		return stderrs.Internal.
			EmbedErrors(err).
			SetMessage(err.Error())
	}

	return nil
}

func (c *Cache[K, V]) Get(ctx context.Context, key K) (V, *stderrs.Error) {
	var value V

	err := c.storage.Get(ctx, fmt.Sprintf("%v", key), &value)
	if err != nil {
		return value, stderrs.Internal.
			EmbedErrors(err).
			SetMessage(err.Error())
	}

	return value, nil
}

func (c *Cache[K, V]) Delete(ctx context.Context, keys ...K) *stderrs.Error {
	for _, key := range keys {
		err := c.storage.Delete(ctx, fmt.Sprintf("%v", key))
		if err != nil {
			return stderrs.Internal.
				EmbedErrors(err).
				SetMessage(err.Error())
		}
	}

	return nil
}

func (c *Cache[K, V]) GC() *stderrs.Error {
	return nil
}

package redis

import (
	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Redis         redis.RingOptions
	LocalCache    cache.LocalCache
	MarshalFunc   cache.MarshalFunc
	UnmarshalFunc cache.UnmarshalFunc
}

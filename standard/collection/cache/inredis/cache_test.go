package inredis_test

import (
	"context"
	"testing"

	"github.com/auvitly/go-tools/standard/collection/cache"
	"github.com/auvitly/go-tools/standard/collection/cache/inredis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	c, stderr := inredis.New[int, int](inredis.Config{
		Redis: redis.RingOptions{
			Addrs: map[string]string{
				"server1": ":6379",
			},
		},
	})
	assert.Nil(t, stderr)

	stderr = c.Set(context.Background(), 1, 1, cache.Options{})
	assert.Nil(t, stderr)

	value, stderr := c.Get(context.Background(), 1)
	assert.Nil(t, stderr)

	assert.Equal(t, value, 1)
}

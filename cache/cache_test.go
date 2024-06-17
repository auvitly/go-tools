package cache_test

import (
	"context"
	"github.com/auvitly/go-tools/cache"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGetTTL(t *testing.T) {
	var c cache.Cache[string, string]

	c.Set("key", "value", cache.WithTTL(time.Second))

	_, ok := c.Get("key")
	require.True(t, ok)

	time.Sleep(2 * time.Second)

	_, ok = c.Get("key")
	require.False(t, ok)
}

func TestGetContextDeadline(t *testing.T) {
	var c cache.Cache[string, string]

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	c.Set("key", "value", cache.WithContext(ctx))

	_, ok := c.Get("key")
	require.True(t, ok)

	time.Sleep(2 * time.Second)

	_, ok = c.Get("key")
	require.False(t, ok)
}

func TestDelete(t *testing.T) {
	var c cache.Cache[string, string]

	c.Set("key", "value")

	_, ok := c.Get("key")
	require.True(t, ok)

	time.Sleep(2 * time.Second)

	c.Delete("key")
	_, ok = c.Get("key")
	require.False(t, ok)
}

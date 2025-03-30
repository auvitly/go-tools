package inmemory_test

import (
	"context"
	"testing"

	"github.com/auvitly/go-tools/standard/collection/cache"
	"github.com/auvitly/go-tools/standard/collection/cache/inmemory"
)

func TestCache(t *testing.T) {
	var c cache.Cache[int, int] = inmemory.New[int, int](inmemory.Config{
		DefaultTTL:  0,
		RecordLimit: 100,
	})

	for i := 0; i < 500; i++ {
		c.Set(context.Background(), i, i, cache.Options{})
	}

	c.GC()

	t.Log(c)
}

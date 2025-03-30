package inmemory_test

import (
	"context"
	"testing"

	"github.com/auvitly/go-tools/standard/collection/cache"
	"github.com/auvitly/go-tools/standard/collection/cache/inmemory"
)

func TestCache(t *testing.T) {
	var c cache.Cache[int, int] = inmemory.New[int, int](inmemory.Config{
		RecordLifeTime: 0,
		RecordLimit:    100,
	})

	for i := 0; i < 500; i++ {
		c.Set(context.Background(), i, cache.Item[int]{Value: i})
	}

	c.GC(context.Background())

	t.Log(c)
}

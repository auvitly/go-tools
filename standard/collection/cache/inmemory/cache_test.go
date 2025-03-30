package inmemory_test

import (
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
		c.Set(i, cache.Item[int]{Value: i})
	}

	c.GC()

	t.Log(c)
}

// Package cache provides an in memory cache model.
package cache

import (
	"context"
	"time"
)

// Cache - unified interface model.
type Cache[K comparable, V any] interface {
	Get(ctx context.Context, key K) Item[V]
	Lookup(ctx context.Context, key K) (Item[V], bool)
	Set(ctx context.Context, key K, item Item[V])
	Delete(ctx context.Context, keys ...K)
	GC(ctx context.Context)
}

// Item - unificated item model.
type Item[V any] struct {
	Value    V
	Deadline *time.Time
}

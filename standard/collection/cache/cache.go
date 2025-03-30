// Package cache provides an in memory cache model.
package cache

import (
	"context"
	"time"

	"github.com/auvitly/go-tools/stderrs"
)

// Cache - unified interface model.
type Cache[K comparable, V any] interface {
	Get(key K) (Item[V], *stderrs.Error)
	GetWithContext(ctx context.Context, key K) (Item[V], *stderrs.Error)
	Set(key K, item Item[V]) *stderrs.Error
	SetWithContext(ctx context.Context, key K, item Item[V]) *stderrs.Error
	Delete(keys ...K) *stderrs.Error
	DeleteWithContext(ctx context.Context, keys ...K) *stderrs.Error
	GC() *stderrs.Error
}

// Item - unificated item model.
type Item[V any] struct {
	Value    V
	Deadline *time.Time
}

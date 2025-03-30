// Package cache provides an in memory cache model.
package cache

import (
	"context"
	"time"

	"github.com/auvitly/go-tools/stderrs"
)

// Cache - unified interface model.
type Cache[K comparable, R any] interface {
	Get(ctx context.Context, key K) (R, *stderrs.Error)
	Set(ctx context.Context, key K, item R, options Options) *stderrs.Error
	Delete(ctx context.Context, keys ...K) *stderrs.Error
	GC() *stderrs.Error
}

type Options struct {
	TTL *time.Duration
}

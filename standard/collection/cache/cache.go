// Package cache provides an in memory cache model.
package cache

import (
	"time"
)

// Cache - unified interface model.
type Cache[K comparable, V any] interface {
	Get(key K) Item[V]
	Lookup(key K) (Item[V], bool)
	Set(key K, item Item[V])
	Delete(keys ...K)
	GC()
}

// Item - unificated item model.
type Item[V any] struct {
	Value    V
	Deadline *time.Time
}

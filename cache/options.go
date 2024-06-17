package cache

import (
	"context"
	"time"
)

// Option - methods for tuning cache behavior.
type Option func(config *config)

// WithTTL - set record lifetime.
func WithTTL(ttl time.Duration) Option {
	return func(config *config) {
		config.TTL = ttl
	}
}

// WithTimestamp - set record lifetime.
func WithTimestamp(ts time.Time) Option {
	return func(config *config) {
		config.TTL = time.Until(ts)
	}
}

// WithContext - set the lifetime of a record by context.
func WithContext(ctx context.Context) Option {
	return func(config *config) {
		config.Context = ctx
	}
}

package cache

import (
	"context"
	"github.com/auvitly/go-tools/cache/internal"
	"time"
)

type Config struct {
	TTL     time.Duration
	Context context.Context
}

func (c Config) getExpirations() (expirations []internal.Expiration) {
	if c.TTL != 0 {
		expirations = append(expirations, &internal.ExpirationTTL{TS: time.Now().Add(c.TTL)})
	}

	if c.Context != nil {
		expirations = append(expirations, &internal.ExpirationContext{Context: c.Context})
	}

	return expirations
}

type Option func(config *Config)

// WithTTL - set record lifetime.
func WithTTL(ttl time.Duration) Option {
	return func(config *Config) {
		config.TTL = ttl
	}
}

// WithTimestamp - set record lifetime.
func WithTimestamp(ts time.Time) Option {
	return func(config *Config) {
		config.TTL = time.Until(ts)
	}
}

// WithContext - set the lifetime of a record by context.
func WithContext(ctx context.Context) Option {
	return func(config *Config) {
		config.Context = ctx
	}
}

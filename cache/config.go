package cache

import (
	"context"
	"github.com/auvitly/go-tools/cache/internal"
	"time"
)

type config struct {
	TTL     time.Duration
	Context context.Context
}

func (c config) getExpirations() (expirations []internal.Expiration) {
	if c.TTL != 0 {
		expirations = append(expirations, &internal.ExpirationTTL{TS: time.Now().Add(c.TTL)})
	}

	if c.Context != nil {
		expirations = append(expirations, &internal.ExpirationContext{Context: c.Context})
	}

	return expirations
}

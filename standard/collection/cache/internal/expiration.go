package internal

import (
	"context"
	"time"
)

type Expiration interface {
	IsExpired() bool
}

type ExpirationTTL struct {
	TS time.Time
}

func (e *ExpirationTTL) IsExpired() bool {
	if time.Until(e.TS) < 0 {
		return true
	}

	return false
}

type ExpirationContext struct {
	Context context.Context
}

func (e *ExpirationContext) IsExpired() bool {
	if e.Context == nil {
		return true
	}

	if deadline, ok := e.Context.Deadline(); ok && time.Until(deadline) < 0 {
		return true
	}

	return false
}

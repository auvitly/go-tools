package inmemory

import (
	"time"
)

type Config struct {
	DefaultTTL  time.Duration
	RecordLimit int
}

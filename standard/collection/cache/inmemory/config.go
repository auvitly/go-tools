package inmemory

import (
	"time"
)

type Config struct {
	RecordLifeTime time.Duration
	RecordLimit    int
}

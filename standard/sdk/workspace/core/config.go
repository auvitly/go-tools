package core

import "time"

type Config struct {
	TaskDowntime    time.Duration
	PullingInterval time.Duration
}

package core

import "time"

type Config struct {
	SessionDecayTime time.Duration
	PullingInterval  time.Duration
}

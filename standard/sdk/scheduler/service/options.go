package service

import "time"

type Option func(*config)

type config struct {
	UnassignTaskInactionInterval time.Duration
	DeleteCompletedSessions      bool
	DeleteCompletedTasks         bool
}

func WithUnassignTaskInactionInterval(interval time.Duration) func(c *config) {
	return func(cfg *config) {
		cfg.UnassignTaskInactionInterval = interval
	}
}

func WithDeleteCompletedTasks(cfg *config) {
	cfg.DeleteCompletedTasks = true
}

func WithDeleteCompletedSessions(cfg *config) {
	cfg.DeleteCompletedSessions = true
}

package service

import (
	"github.com/auvitly/go-tools/collection/stderrs"
	"github.com/auvitly/go-tools/standard/sdk/scheduler/entity/session"
	"github.com/auvitly/go-tools/standard/sdk/scheduler/entity/task"
	"github.com/auvitly/go-tools/standard/sdk/scheduler/entity/worker"
)

type Service[T task.IsTask, W worker.IsWorker, S session.IsSession] struct {
	config       config
	dependencies Dependencies[T, W, S]
}

type Dependencies[T task.IsTask, W worker.IsWorker, S session.IsSession] struct {
	TaskStorage    task.Storage[T]
	WorkerStorage  worker.Storage[W]
	SessionStorage session.Storage[S]
}

func New[
	T task.IsTask,
	W worker.IsWorker,
	S session.IsSession,
](
	dps Dependencies[T, W, S],
	options ...Option,
) (
	*Service[T, W, S],
	*stderrs.Error,
) {
	var cfg config

	for _, option := range options {
		option(&cfg)
	}

	return &Service[T, W, S]{
		config:       cfg,
		dependencies: dps,
	}, nil
}

package service

import (
	"cmp"

	"github.com/auvitly/go-tools/standard/sdk/scheduler/storage"
)

type Dependencies[T, M, S cmp.Ordered] struct {
	WorkerStorage  storage.WorkerStorage[T]
	SessionStorage storage.SessionStorage
	TaskStorage    storage.TaskStorage[T, M, S]
}

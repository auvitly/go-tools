package core

import (
	"cmp"

	"github.com/auvitly/go-tools/standard/workspace/storage"
)

type Dependencies[T, M cmp.Ordered] struct {
	TaskStorage    storage.TaskStorage[T, M]
	SessionStorage storage.SessionStorage[T]
	WorkerStorage  storage.WorkerStorage[T]
}

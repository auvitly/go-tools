package inmemory

import (
	"context"
	"sync"

	"github.com/auvitly/go-tools/collection/stderrs"
	"github.com/auvitly/go-tools/standard/sdk/scheduler/entity/task"
	"github.com/google/uuid"
)

type TaskStorage[T task.IsTask] struct {
	mu      sync.RWMutex
	storage map[uuid.UUID]T
}

func NewTaskStorage[T task.IsTask]() *TaskStorage[T] {
	return &TaskStorage[T]{
		storage: make(map[uuid.UUID]T),
	}
}

func (s *TaskStorage[T]) Get(ctx context.Context, id uuid.UUID) (T, *stderrs.Error) {

}

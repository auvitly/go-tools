package inmemory

import (
	"context"
	"sync"

	"github.com/auvitly/go-tools/collection/stderrs"
	"github.com/auvitly/go-tools/standard/sdk/scheduler/entity/task"
	"github.com/auvitly/go-tools/standard/utils/reflector"
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
	s.mu.RLock()
	defer s.mu.Unlock()

	record, ok := s.storage[id]
	if !ok {
		return reflector.Nil[T](), stderrs.NotFound.SetMessage("not found task with id '%s'", id.String())
	}

	return record, nil
}

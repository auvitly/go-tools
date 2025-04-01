package inmemory

import (
	"cmp"
	"context"
	"sync"
	"time"

	"github.com/auvitly/go-tools/standard/sdk/workspace/entity"
	"github.com/auvitly/go-tools/standard/sdk/workspace/storage"
	"github.com/auvitly/go-tools/stderrs"
	"github.com/google/uuid"
)

type WorkerStorage[T cmp.Ordered] struct {
	mu      sync.RWMutex
	storage map[uuid.UUID]*entity.Worker[T]
}

func NewWorkerStorage[T cmp.Ordered]() *WorkerStorage[T] {
	return &WorkerStorage[T]{
		storage: map[uuid.UUID]*entity.Worker[T]{},
	}
}

func (s *WorkerStorage[T]) Save(ctx context.Context, params storage.WorkerSaveParams[T]) (*entity.Worker[T], *stderrs.Error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	worker, ok := s.storage[params.WorkerID]
	if !ok {
		s.storage[params.WorkerID] = &entity.Worker[T]{
			ID:        params.WorkerID,
			Type:      params.Type,
			Version:   params.Version,
			Labels:    params.Labels,
			CreatedAT: time.Now(),
			UpdatedAT: time.Now(),
		}
	} else {
		worker.Type = params.Type
		worker.Version = params.Version
		worker.UpdatedAT = time.Now()
	}

	return s.storage[params.WorkerID], nil
}

func (s *WorkerStorage[T]) Get(ctx context.Context, params storage.WorkerGetParams) (*entity.Worker[T], *stderrs.Error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	worker, ok := s.storage[params.WorkerID]
	if !ok {
		return nil, stderrs.NotFound.SetMessage("not found worker with id=%s", params.WorkerID.String())
	}

	return worker, nil
}

func (s *WorkerStorage[T]) Delete(ctx context.Context, params storage.WorkerDeleteParams) *stderrs.Error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.storage, params.WorkerID)

	return nil
}

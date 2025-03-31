package inmemory

import (
	"cmp"
	"context"
	"sync"
	"time"

	"github.com/auvitly/go-tools/standard/workspace/entity"
	"github.com/auvitly/go-tools/standard/workspace/storage"
	"github.com/auvitly/go-tools/stderrs"
	"github.com/google/uuid"
)

type TaskStorage[T, M, S cmp.Ordered] struct {
	mu      sync.RWMutex
	storage map[uuid.UUID]*entity.Task[T, M, S]
}

func NewTaskStorage[T, M, S cmp.Ordered]() *TaskStorage[T, M, S] {
	return &TaskStorage[T, M, S]{
		storage: map[uuid.UUID]*entity.Task[T, M, S]{},
	}
}

func (s *TaskStorage[T, M, S]) Update(ctx context.Context, params storage.TaskUpdateParams[S]) (*entity.Task[T, M, S], *stderrs.Error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.storage[params.TaskID]
	if !ok {
		return nil, stderrs.NotFound.SetMessage("not found task with id=%s", params.TaskID.String())
	}

	task.StatusCode = params.StatusCode
	task.StateData = params.StateData
	task.Result = params.Result
	task.UpdatedTS = params.UpdatedAT
	task.CatchLaterTS = params.CatchLaterAT
	task.DoneTS = params.DoneTS
	task.SessionID = params.SessionID
	task.AssignTS = params.AssignTS

	return task, nil
}

func (s *TaskStorage[T, M, S]) Push(ctx context.Context, params storage.TaskPushParams[T, M, S]) (*entity.Task[T, M, S], *stderrs.Error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var (
		id = uuid.New()
		ts = time.Now()
	)

	s.storage[id] = &entity.Task[T, M, S]{
		ID:           id,
		ParentTaskID: params.ParentTaskID,
		Type:         params.Type,
		Mode:         params.Mode,
		StatusCode:   nil,
		Args:         params.Args,
		StateData:    nil,
		Result:       nil,
		CreatedTS:    ts,
		UpdatedTS:    ts,
		CatchLaterTS: nil,
		DoneTS:       nil,
		SessionID:    nil,
		AssignTS:     nil,
		Labels:       params.Labels,
	}

	return s.storage[id], nil
}

func (s *TaskStorage[T, M, S]) Pop(ctx context.Context, params storage.TaskPopParams[T]) (*entity.Task[T, M, S], *stderrs.Error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var ts = time.Now()

loop:
	for _, task := range s.storage {
		if task.SessionID != nil || task.AssignTS != nil || task.DoneTS != nil {
			continue loop
		}

		if task.Type != params.Type {
			continue loop
		}

		for key, value := range params.Labels {
			label, ok := task.Labels[key]
			if !ok {
				continue
			}

			if label != value {
				continue loop
			}
		}

		task.AssignTS = &ts
		task.SessionID = &params.SessionID

		return task, nil
	}

	return nil, stderrs.NotFound.SetMessage("not found task for worker")
}

func (s *TaskStorage[T, M, S]) Get(ctx context.Context, params storage.TaskGetParams) (*entity.Task[T, M, S], *stderrs.Error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, ok := s.storage[params.TaskID]
	if !ok {
		return nil, stderrs.NotFound.SetMessage("not found task with id=%s", params.TaskID.String())
	}

	return task, nil
}

func (s *TaskStorage[T, M, S]) Flush(ctx context.Context, params storage.TaskFlushParams) *stderrs.Error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, task := range s.storage {
		if task.AssignTS == nil || task.SessionID == nil {
			continue
		}

		task.SessionID = nil
		task.AssignTS = nil
	}

	return nil
}

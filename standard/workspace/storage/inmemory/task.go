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

type TaskStorage[T, M cmp.Ordered] struct {
	mu      sync.RWMutex
	storage map[uuid.UUID]*entity.Task[T, M]
}

func NewTaskStorage[T, M cmp.Ordered]() *TaskStorage[T, M] {
	return &TaskStorage[T, M]{
		storage: map[uuid.UUID]*entity.Task[T, M]{},
	}
}

func (s *TaskStorage[T, M]) Update(ctx context.Context, params storage.TaskUpdateParams) (*entity.Task[T, M], *stderrs.Error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.storage[params.TaskID]
	if !ok {
		return nil, stderrs.NotFound.SetMessage("not found task with id=%s", params.TaskID.String())
	}

	task.Status = params.Status
	task.State = params.State
	task.Result = params.Result
	task.UpdatedTS = params.UpdatedAT
	task.CatchLaterTS = params.CatchLaterAT
	task.DoneTS = params.DoneTS
	task.SessionID = params.SessionID
	task.AssignTS = params.AssignTS

	return task, nil
}

func (s *TaskStorage[T, M]) Push(ctx context.Context, params storage.TaskPushParams[T, M]) (*entity.Task[T, M], *stderrs.Error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var (
		id = uuid.New()
		ts = time.Now()
	)

	s.storage[id] = &entity.Task[T, M]{
		ID:           id,
		ParentTaskID: params.ParentTaskID,
		Type:         params.Type,
		Mode:         params.Mode,
		Status:       params.Status,
		Args:         params.Args,
		State:        nil,
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

func (s *TaskStorage[T, M]) Pop(ctx context.Context, params storage.TaskPopParams[T]) (*entity.Task[T, M], *stderrs.Error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var ts = time.Now()

loop:
	for _, task := range s.storage {
		if task.SessionID != nil || task.AssignTS != nil {
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

func (s *TaskStorage[T, M]) Get(ctx context.Context, params storage.TaskGetParams) (*entity.Task[T, M], *stderrs.Error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, ok := s.storage[params.TaskID]
	if !ok {
		return nil, stderrs.NotFound.SetMessage("not found task with id=%s", params.TaskID.String())
	}

	return task, nil
}

func (s *TaskStorage[T, M]) Flush(ctx context.Context, params storage.TaskFlushParams) *stderrs.Error {
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

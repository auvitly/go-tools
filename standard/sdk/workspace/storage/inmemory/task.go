package inmemory

import (
	"cmp"
	"context"
	"maps"
	"slices"
	"sync"
	"time"

	"github.com/auvitly/go-tools/standard/sdk/workspace/entity"
	"github.com/auvitly/go-tools/standard/sdk/workspace/storage"
	"github.com/auvitly/go-tools/stderrs"
	"github.com/google/uuid"
)

type TaskStorage[T, M, S cmp.Ordered] struct {
	mu      sync.RWMutex
	storage map[uuid.UUID]*entity.Task[T, M, S]
	config  TaskConfig
}

type TaskConfig struct {
	DeleteCompleted bool
}

func NewTaskStorage[T, M, S cmp.Ordered](config TaskConfig) *TaskStorage[T, M, S] {
	return &TaskStorage[T, M, S]{
		storage: map[uuid.UUID]*entity.Task[T, M, S]{},
		config:  config,
	}
}

func cloneTask[T, M, S cmp.Ordered](t *entity.Task[T, M, S]) *entity.Task[T, M, S] {
	var task = *t

	task.Args = slices.Clone(task.Args)
	task.SessionID = clonePtr(task.SessionID)
	task.CatchLaterTS = clonePtr(task.CatchLaterTS)
	task.ParentTaskID = clonePtr(task.ParentTaskID)
	task.StatusCode = clonePtr(task.StatusCode)
	task.StateData = slices.Clone(task.StateData)
	task.DoneTS = clonePtr(task.DoneTS)
	task.AssignTS = clonePtr(task.AssignTS)
	task.Labels = maps.Clone(task.Labels)

	return &task
}

func clonePtr[T any](value *T) *T {
	if value == nil {
		return nil
	}

	var cloned = *value

	return &cloned
}

func (s *TaskStorage[T, M, S]) Update(ctx context.Context, params storage.TaskUpdateParams[S]) (*entity.Task[T, M, S], *stderrs.Error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.storage[params.TaskID]
	if !ok {
		return nil, stderrs.NotFound.SetMessage("not found task with id=%s", params.TaskID.String())
	}

	if task.DoneTS != nil {
		return nil, stderrs.Aborted.SetMessage("task already done")
	}

	if s.config.DeleteCompleted && params.DoneTS != nil {
		defer delete(s.storage, params.TaskID)
	}

	task.StatusCode = params.StatusCode
	task.StateData = params.StateData
	task.Result = params.Result
	task.UpdatedTS = params.UpdatedTS
	task.CatchLaterTS = params.CatchLaterTS
	task.DoneTS = params.DoneTS
	task.SessionID = params.SessionID
	task.AssignTS = params.AssignTS

	return cloneTask(task), nil
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

	return cloneTask(s.storage[id]), nil
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

		return cloneTask(task), nil
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

	return cloneTask(task), nil
}

func (s *TaskStorage[T, M, S]) List(ctx context.Context, params storage.TaskListParams) ([]*entity.Task[T, M, S], *stderrs.Error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var tasks []*entity.Task[T, M, S]

	for _, task := range s.storage {
		if task.AssignTS != nil && params.OnlyAssigned {
			tasks = append(tasks, cloneTask(task))
		}
	}

	return tasks, nil
}

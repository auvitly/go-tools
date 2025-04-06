package inmemory

import (
	"context"
	"maps"
	"slices"
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

	return record.Clone().(T), nil
}

func (s *TaskStorage[T]) List(ctx context.Context, filters ...task.IsListFilter) ([]T, *stderrs.Error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stderr := supportedTaskFilterList(filters...)
	if stderr != nil {
		return nil, stderr
	}

	var (
		storage = maps.Clone(s.storage)
		list    []T
	)

	for _, item := range storage {
		if checkTaskFilterList(item, filters...) {
			list = append(list, item.Clone().(T))
		}
	}

	return list, nil
}

func supportedTaskFilterList(filters ...task.IsListFilter) *stderrs.Error {
	for _, filter := range filters {
		switch value := filter.(type) {
		case task.ListFilterTypes:
		case task.ListFilterLabels:
		case task.ListFilterModes:
		case task.ListFilterStatuses:
		case task.ListFilterOnlyAssigned:
		case task.ListFilterWithoutAssigned:
		case task.ListFilterPagination:
		case task.ListFilterSortByCreatedAT:
		default:
			return stderrs.Unimplemented.SetMessage("unimplemented filter '%T'", value)
		}
	}

	return nil
}

func checkTaskFilterList[T task.IsTask](item T, filters ...task.IsListFilter) bool {
	for _, filter := range filters {
		switch t := filter.(type) {
		case task.ListFilterTypes:
			if !slices.Contains(t.Types, item.Impl().Type) {
				return false
			}
		case task.ListFilterLabels:
			for key, label := range t.WorkerLabels {
				if val, exists := item.Impl().WorkerLabels[key]; !exists || val != label {
					return false
				}
			}
		case task.ListFilterModes:
			if !slices.Contains(t.Modes, item.Impl().Mode) {
				return false
			}
		case task.ListFilterStatuses:
			if !slices.Contains(t.Statuses, item.Impl().Status) {
				return false
			}
		case task.ListFilterOnlyAssigned:
			if item.Impl().WorkerAssignTS == nil && item.Impl().WorkerSessionID == nil {
				return false
			}
		case task.ListFilterWithoutAssigned:
			if item.Impl().WorkerAssignTS != nil && item.Impl().WorkerSessionID != nil {
				return false
			}
		}
	}

	return true
}

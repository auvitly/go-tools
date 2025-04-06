package task

import (
	"context"
	"time"

	"github.com/auvitly/go-tools/collection/stderrs"
	"github.com/google/uuid"
)

type Storage[T IsTask] interface {
	Get(ctx context.Context, taskID uuid.UUID) (T, *stderrs.Error)
	List(ctx context.Context, filters ...IsListFilter) ([]T, *stderrs.Error)
	Push(ctx context.Context, task IsTask) (T, *stderrs.Error)
}

type PopParams struct {
	SessionID uuid.UUID
	Types     []string
	Modes     []Mode
	Labels    map[string]string
}

type IsListFilter interface{ implListFilter() ListFilter }

type ListFilter struct{}

type (
	ListFilterIDs struct {
		ListFilter
		IDs []uuid.UUID
	}
	ListFilterParentIDs struct {
		ListFilter
		ParentIDs []uuid.UUID
	}
	ListFilterTypes struct {
		ListFilter
		Types []string
	}
	ListFilterModes struct {
		ListFilter
		Modes []Mode
	}
	ListFilterStatuses struct {
		ListFilter
		Statuses []Status
	}

	ListFilterInactive struct {
		ListFilter
		Interval time.Duration
	}
	ListFilterOnlyAssigned struct {
		ListFilter
	}
	ListFilterWithoutAssigned struct {
		ListFilter
	}
	ListFilterLabels struct {
		ListFilter
		WorkerLabels map[string]string
	}
	ListFilterPagination struct {
		ListFilter
		Limit  int
		Offset int
	}
)

func (f ListFilter) implListFilter() ListFilter { return f }

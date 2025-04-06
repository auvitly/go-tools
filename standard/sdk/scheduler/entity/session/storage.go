package session

import (
	"context"

	"github.com/auvitly/go-tools/collection/stderrs"
	"github.com/google/uuid"
)

type Storage[S IsSession] interface {
	Receive(ctx context.Context, workerID uuid.UUID) (S, *stderrs.Error)
	Get(ctx context.Context, sessionID uuid.UUID) (S, *stderrs.Error)
	List(ctx context.Context, filters ...IsListFilter) ([]S, *stderrs.Error)
	Save(ctx context.Context, session S) (S, *stderrs.Error)
	Drop(ctx context.Context, sessionID uuid.UUID) *stderrs.Error
}

type IsListFilter interface{ implListFilter() ListFilter }

type ListFilter struct{}

type (
	ListFilterIDs struct {
		ListFilter
		IDs []uuid.UUID
	}
	ListFilterWorkerIDs struct {
		ListFilter
		WorkerIDs []uuid.UUID
	}
	ListFilterOnlyCompleted struct {
		ListFilter
	}
	ListFilterOnlyUncompleted struct {
		ListFilter
	}
)

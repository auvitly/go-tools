package storage

import (
	"context"

	"github.com/auvitly/go-tools/standard/workspace/entity"
	"github.com/auvitly/go-tools/stderrs"
	"github.com/google/uuid"
)

type SessionStorage interface {
	New(ctx context.Context, params SessionNewParams) (*entity.Session, *stderrs.Error)
	Get(ctx context.Context, params SessionGetParams) (*entity.Session, *stderrs.Error)
	List(ctx context.Context, params SessionListParams) ([]*entity.Session, *stderrs.Error)
	Drop(ctx context.Context, params SessionDropParams) *stderrs.Error
	Done(ctx context.Context, params SessionDoneParams) (*entity.Session, *stderrs.Error)
}

type SessionNewParams struct {
	WorkerID uuid.UUID
}

type SessionGetParams struct {
	SessionID uuid.UUID
}

type SessionListParams struct {
	WorkerID      *uuid.UUID
	SessionIDs    []uuid.UUID
	OnlyCompleted bool
}

type SessionDropParams struct {
	SessionID uuid.UUID
}

type SessionDoneParams struct {
	SessionID uuid.UUID
}

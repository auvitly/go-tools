package storage

import (
	"cmp"
	"context"
	"time"

	"github.com/auvitly/go-tools/standard/sdk/workspace/entity"
	"github.com/auvitly/go-tools/stderrs"
	"github.com/google/uuid"
)

type WorkerStorage[T cmp.Ordered] interface {
	Save(ctx context.Context, params WorkerSaveParams[T]) (*entity.Worker[T], *stderrs.Error)
	Get(ctx context.Context, params WorkerGetParams) (*entity.Worker[T], *stderrs.Error)
	Delete(ctx context.Context, params WorkerDeleteParams) *stderrs.Error
}

type WorkerSaveParams[T cmp.Ordered] struct {
	WorkerID  uuid.UUID
	Type      T
	Version   string
	Labels    map[string]string
	CreatedAT time.Time
	UpdatedAT time.Time
}

type WorkerGetParams struct {
	WorkerID uuid.UUID
}

type WorkerDeleteParams struct {
	WorkerID uuid.UUID
}

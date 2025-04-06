package worker

import (
	"context"

	"github.com/auvitly/go-tools/collection/stderrs"
	"github.com/google/uuid"
)

type Storage[W IsWorker] interface {
	Save(ctx context.Context, worker W) (W, *stderrs.Error)
	Get(ctx context.Context, workerID uuid.UUID) (W, *stderrs.Error)
	Delete(ctx context.Context, workerID uuid.UUID) *stderrs.Error
}

package worker

import (
	"time"

	"github.com/google/uuid"
)

type Worker struct {
	ID        uuid.UUID
	Types     []string
	Labels    map[string]string
	CreatedAT time.Time
	UpdatedAT time.Time
}

type IsWorker interface{ Impl() *Worker }

func (w *Worker) Impl() *Worker { return w }

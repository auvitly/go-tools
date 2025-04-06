package task

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID       uuid.UUID
	ParentID *uuid.UUID

	Type   string
	Mode   Mode
	Status Status

	Argument map[string]any
	Data     map[string]any
	Result   map[string]any

	CreatedTS time.Time
	UpdatedTS time.Time
	PendingTS time.Time

	WorkerSessionID *uuid.UUID
	WorkerAssignTS  *time.Time
	WorkerLabels    map[string]string
}

type IsTask interface{ Impl() *Task }

type (
	Status string
	Mode   string
)

const (
	StatusCreated      Status = "created"
	StatusInProgress   Status = "in_progress"
	StatusError        Status = "error"
	StatusCompensating Status = "compensating"
	StatusCanceled     Status = "canceled"
	StatusDone         Status = "done"
	StatusCompleted    Status = "completed"
)

const (
	ModeGlobalSync  Mode = "global_sync"
	ModeLocalSync   Mode = "local_sync"
	ModeGlobalAsync Mode = "global_async"
	ModeLocalAsync  Mode = "local_async"
	ModeManual      Mode = "manual"
)

func (t *Task) Impl() *Task { return t }

package task

import (
	"slices"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID       uuid.UUID
	ParentID *uuid.UUID

	Type   string
	Mode   Mode
	Status Status

	Arguments map[string]any
	Data      map[string]any
	Results   map[string]any

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
	StatusPending      Status = "pending"
	StatusCreated      Status = "created"
	StatusInProgress   Status = "in_progress"
	StatusError        Status = "error"
	StatusException    Status = "exception"
	StatusCompensating Status = "compensating"
	StatusCanceled     Status = "canceled"
	StatusDone         Status = "done"
	StatusCompleted    Status = "completed"
)

var Statuses = []Status{
	StatusPending,
	StatusCreated,
	StatusInProgress,
	StatusError,
	StatusException,
	StatusCompensating,
	StatusCanceled,
	StatusDone,
	StatusCompleted,
}

func (s Status) Valid() bool {
	return slices.Contains(Statuses, s)
}

const (
	ModeAsync  Mode = "async"
	ModeSync   Mode = "sync"
	ModeManual Mode = "manual"
)

func (t *Task) Impl() *Task { return t }

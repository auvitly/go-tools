package storage

import (
	"cmp"
	"context"
	"encoding/json"
	"time"

	"github.com/auvitly/go-tools/standard/workspace/entity"
	"github.com/auvitly/go-tools/stderrs"
	"github.com/google/uuid"
)

type TaskStorage[T, M cmp.Ordered] interface {
	Update(ctx context.Context, params TaskUpdateParams) (*entity.Task[T, M], *stderrs.Error)
	Push(ctx context.Context, params TaskPushParams[T, M]) (*entity.Task[T, M], *stderrs.Error)
	Pop(ctx context.Context, params TaskPopParams[T]) (*entity.Task[T, M], *stderrs.Error)
	Get(ctx context.Context, params TaskGetParams) (*entity.Task[T, M], *stderrs.Error)
	Flush(ctx context.Context, params TaskFlushParams) *stderrs.Error
}

type TaskUpdateParams struct {
	TaskID       uuid.UUID
	Status       string
	State        json.RawMessage
	Result       *json.RawMessage
	UpdatedAT    time.Time
	CatchLaterAT *time.Time
	DoneTS       *time.Time
	SessionID    *uuid.UUID
	AssignTS     *time.Time
}

type TaskPushParams[T, M cmp.Ordered] struct {
	ParentTaskID *uuid.UUID
	Type         T
	Mode         M
	Status       string
	Args         json.RawMessage
	Labels       map[string]string
}

type TaskPopParams[T cmp.Ordered] struct {
	SessionID uuid.UUID
	Type      T
	Labels    map[string]string
}

type TaskGetParams struct {
	TaskID uuid.UUID
}

type TaskFlushParams struct {
	Downtime time.Duration
}

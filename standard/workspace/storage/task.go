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

type TaskStorage[T, M, S cmp.Ordered] interface {
	Update(ctx context.Context, params TaskUpdateParams[S]) (*entity.Task[T, M, S], *stderrs.Error)
	Push(ctx context.Context, params TaskPushParams[T, M, S]) (*entity.Task[T, M, S], *stderrs.Error)
	Pop(ctx context.Context, params TaskPopParams[T]) (*entity.Task[T, M, S], *stderrs.Error)
	Get(ctx context.Context, params TaskGetParams) (*entity.Task[T, M, S], *stderrs.Error)
	Flush(ctx context.Context, params TaskFlushParams) *stderrs.Error
}

type TaskUpdateParams[S cmp.Ordered] struct {
	TaskID       uuid.UUID
	StatusCode   S
	State        json.RawMessage
	Result       *json.RawMessage
	UpdatedAT    time.Time
	CatchLaterAT *time.Time
	DoneTS       *time.Time
	SessionID    *uuid.UUID
	AssignTS     *time.Time
}

type TaskPushParams[T, M, S cmp.Ordered] struct {
	ParentTaskID *uuid.UUID
	Type         T
	Mode         M
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

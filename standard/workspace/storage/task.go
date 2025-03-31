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
	Untie(ctx context.Context, params TaskUntieParams) *stderrs.Error
}

type TaskUpdateParams struct {
	ID           uuid.UUID
	SessionID    uuid.UUID
	Status       string
	State        json.RawMessage
	Result       *json.RawMessage
	CatchLaterAT *time.Time
	Labels       map[string]string
}

type TaskPushParams[T, M cmp.Ordered] struct {
	ParentTaskID *uuid.UUID
	Type         T
	Mode         M
	In           json.RawMessage
	Labels       map[string]string
}

type TaskPopParams[T cmp.Ordered] struct {
	SessionID uuid.UUID
	Type      T
	Labels    map[string]string
}

type TaskUntieParams struct {
	Downtime time.Duration
}

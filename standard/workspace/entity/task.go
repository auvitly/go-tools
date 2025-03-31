package entity

import (
	"cmp"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Task[T, M, S cmp.Ordered] struct {
	ID           uuid.UUID
	ParentTaskID *uuid.UUID
	Type         T
	Mode         M
	Args         json.RawMessage
	StateData    json.RawMessage
	Result       *json.RawMessage
	StatusCode   *S
	CreatedTS    time.Time
	UpdatedTS    time.Time
	CatchLaterTS *time.Time
	DoneTS       *time.Time
	SessionID    *uuid.UUID
	AssignTS     *time.Time
	Labels       map[string]string
}

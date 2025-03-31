package entity

import (
	"cmp"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Task[T, M cmp.Ordered] struct {
	ID           uuid.UUID
	ParentTaskID *uuid.UUID
	Type         T
	Mode         M
	Status       string
	Args         json.RawMessage
	State        json.RawMessage
	Result       *json.RawMessage
	CreatedTS    time.Time
	UpdatedTS    time.Time
	CatchLaterTS *time.Time
	DoneTS       *time.Time
	SessionID    *uuid.UUID
	AssignTS     *time.Time
	Labels       map[string]string
}

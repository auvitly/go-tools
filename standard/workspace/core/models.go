package core

import (
	"cmp"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type CreateTaskParams[T, M cmp.Ordered] struct {
	ParentTaskID *uuid.UUID
	Type         T
	Mode         M
	Status       string
	Args         json.RawMessage
	Labels       map[string]string
}

type ReceiveTaskParams[T cmp.Ordered] struct {
	WorkerID uuid.UUID
	Type     T
	Version  string
	Labels   map[string]string
}

type SetStateParams[T cmp.Ordered] struct {
	TaskID       uuid.UUID
	SessionID    uuid.UUID
	Status       string
	State        json.RawMessage
	Result       *json.RawMessage
	CatchLaterAT *time.Time
}

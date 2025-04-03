package entity

import (
	"cmp"
	"time"

	"github.com/google/uuid"
)

type Data = []byte

type Task[T, M, S cmp.Ordered] SpecificTask[T, M, S, Data, Data]

type SpecificTask[T, M, S cmp.Ordered, A, R any] struct {
	ID           uuid.UUID
	ParentTaskID *uuid.UUID
	Type         T
	Mode         M
	Args         A
	Data         map[string]any
	Result       *R
	StatusCode   *S
	CreatedTS    time.Time
	UpdatedTS    time.Time
	CatchLaterTS *time.Time
	DoneTS       *time.Time
	SessionID    *uuid.UUID
	AssignTS     *time.Time
	Labels       map[string]string
}

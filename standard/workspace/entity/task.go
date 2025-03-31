package entity

import (
	"cmp"
	"encoding/json"
	"maps"
	"slices"
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

func clonePtr[T any](value *T) *T {
	if value == nil {
		return nil
	}

	var cloned = *value

	return &cloned
}

func (t *Task[T, M, S]) Clone() *Task[T, M, S] {
	var task = *t

	task.Args = slices.Clone(task.Args)
	task.SessionID = clonePtr(task.SessionID)
	task.CatchLaterTS = clonePtr(task.CatchLaterTS)
	task.ParentTaskID = clonePtr(task.ParentTaskID)
	task.StatusCode = clonePtr(task.StatusCode)
	task.StateData = slices.Clone(task.StateData)
	task.DoneTS = clonePtr(task.DoneTS)
	task.AssignTS = clonePtr(task.AssignTS)
	task.Labels = maps.Clone(task.Labels)

	return &task
}

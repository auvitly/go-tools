package entity

import (
	"cmp"
	"encoding/json"
	"time"

	"github.com/auvitly/go-tools/stderrs"
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

type SpecifiedTask[T, M, S cmp.Ordered, A, D, R any] struct {
	ID           uuid.UUID
	ParentTaskID *uuid.UUID
	Type         T
	Mode         M
	Args         A
	StateData    D
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

func MarshalTask[T, M, S cmp.Ordered, A, D, R any](task *SpecifiedTask[T, M, S, A, D, R]) (*Task[T, M, S], *stderrs.Error) {
	var (
		args      json.RawMessage
		stateData json.RawMessage
		result    *json.RawMessage
	)

	args, err := json.Marshal(task.Args)
	if err != nil {
		return nil, stderrs.Internal.
			EmbedErrors(err).
			SetMessage("can't marshal args: %s", err.Error())
	}

	stateData, err = json.Marshal(task.StateData)
	if err != nil {
		return nil, stderrs.Internal.
			EmbedErrors(err).
			SetMessage("can't marshal state_data: %s", err.Error())
	}

	if task.Result != nil {
		result = new(json.RawMessage)

		*result, err = json.Marshal(task.Result)
		if err != nil {
			return nil, stderrs.Internal.
				EmbedErrors(err).
				SetMessage("can't marshal results: %s", err.Error())
		}
	}

	return &Task[T, M, S]{
		ID:           task.ID,
		ParentTaskID: task.ParentTaskID,
		Type:         task.Type,
		Mode:         task.Mode,
		Args:         args,
		StateData:    stateData,
		Result:       result,
		StatusCode:   task.StatusCode,
		CreatedTS:    task.CreatedTS,
		UpdatedTS:    task.UpdatedTS,
		CatchLaterTS: task.CatchLaterTS,
		DoneTS:       task.DoneTS,
		SessionID:    task.SessionID,
		AssignTS:     task.AssignTS,
		Labels:       task.Labels,
	}, nil
}

func UnmarshalTask[T, M, S cmp.Ordered, A, D, R any](task *Task[T, M, S]) (*SpecifiedTask[T, M, S, A, D, R], *stderrs.Error) {
	var (
		args      A
		stateData D
		result    *R
	)

	err := json.Unmarshal(task.Args, &args)
	if err != nil {
		return nil, stderrs.Internal.
			EmbedErrors(err).
			SetMessage("can't unmarshal args: %s", err.Error())
	}

	err = json.Unmarshal(task.StateData, &stateData)
	if err != nil {
		return nil, stderrs.Internal.
			EmbedErrors(err).
			SetMessage("can't unmarshal state_data: %s", err.Error())
	}

	err = json.Unmarshal(*task.Result, &result)
	if err != nil {
		return nil, stderrs.Internal.
			EmbedErrors(err).
			SetMessage("can't unmarshal results: %s", err.Error())
	}

	return &SpecifiedTask[T, M, S, A, D, R]{
		ID:           task.ID,
		ParentTaskID: task.ParentTaskID,
		Type:         task.Type,
		Mode:         task.Mode,
		Args:         args,
		StateData:    stateData,
		Result:       result,
		StatusCode:   task.StatusCode,
		CreatedTS:    task.CreatedTS,
		UpdatedTS:    task.UpdatedTS,
		CatchLaterTS: task.CatchLaterTS,
		DoneTS:       task.DoneTS,
		SessionID:    task.SessionID,
		AssignTS:     task.AssignTS,
		Labels:       task.Labels,
	}, nil
}

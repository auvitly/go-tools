package core

import (
	"context"
	"time"

	"github.com/auvitly/go-tools/standard/workspace/entity"
	"github.com/auvitly/go-tools/standard/workspace/storage"
	"github.com/auvitly/go-tools/stderrs"
	"github.com/google/uuid"
)

func (c *Core[T, M, S]) CreateTask(ctx context.Context, params CreateTaskParams[T, M]) (*entity.Task[T, M, S], *stderrs.Error) {
	task, stderr := c.dependencies.TaskStorage.Push(ctx, storage.TaskPushParams[T, M, S]{
		ParentTaskID: params.ParentTaskID,
		Type:         params.Type,
		Mode:         params.Mode,
		Args:         params.Args,
		Labels:       params.Labels,
	})
	if stderr != nil {
		return nil, stderr
	}

	return task, nil
}

func (c *Core[T, M, S]) ReceiveTask(ctx context.Context, params ReceiveTaskParams[T]) (*entity.Task[T, M, S], *stderrs.Error) {
	worker, stderr := c.dependencies.WorkerStorage.Save(ctx, storage.WorkerSaveParams[T]{
		WorkerID:  params.WorkerID,
		Type:      params.Type,
		Version:   params.Version,
		Labels:    params.Labels,
		CreatedAT: time.Now(),
		UpdatedAT: time.Now(),
	})
	if stderr != nil {
		return nil, stderr
	}

	session, stderr := c.dependencies.SessionStorage.New(ctx, storage.SessionNewParams{
		WorkerID: worker.ID,
	})
	if stderr != nil {
		return nil, stderr
	}

	task, stderr := c.dependencies.TaskStorage.Pop(ctx, storage.TaskPopParams[T]{
		SessionID: session.ID,
		Type:      params.Type,
		Labels:    params.Labels,
	})

	switch {
	case stderr.Is(stderrs.NotFound):
		stderr = c.dependencies.SessionStorage.Drop(ctx, storage.SessionDropParams{
			SessionID: session.ID,
		})
		if stderr != nil {
			return nil, stderr
		}

		return nil, stderr
	case stderr != nil:
		return nil, stderr
	default:
		return task, nil
	}
}

func (c *Core[T, M, S]) SetState(ctx context.Context, params SetStateParams[T, S]) *stderrs.Error {
	var ts = time.Now()

	task, stderr := c.dependencies.TaskStorage.Get(ctx, storage.TaskGetParams{
		TaskID: params.TaskID,
	})
	if stderr != nil {
		return stderr
	}

	switch {
	case task.SessionID == nil:
		return stderrs.FailedPrecondition.SetMessage("task unassigned")
	case *task.SessionID != params.SessionID:
		return stderrs.InvalidArgument.SetMessage("worker session not match")
	}

	_, stderr = c.dependencies.TaskStorage.Update(ctx, storage.TaskUpdateParams[S]{
		TaskID: params.TaskID,
		SessionID: func() *uuid.UUID {
			if params.Result != nil || params.CatchLaterAT != nil {
				return nil
			}

			return &params.SessionID
		}(),
		AssignTS: func() *time.Time {
			if params.Result != nil || params.CatchLaterAT != nil {
				return nil
			}

			return &ts
		}(),
		StatusCode:   params.StatusCode,
		State:        params.State,
		Result:       params.Result,
		CatchLaterAT: params.CatchLaterAT,
		UpdatedAT:    ts,
		DoneTS: func() *time.Time {
			if params.Result == nil {
				return nil
			}

			return &ts
		}(),
	})
	if stderr != nil {
		return stderr
	}

	if params.CatchLaterAT == nil && params.Result == nil {
		return nil
	}

	_, stderr = c.dependencies.SessionStorage.Done(ctx, storage.SessionDoneParams{
		SessionID: params.SessionID,
	})
	if stderr != nil {
		return stderr
	}

	return nil
}

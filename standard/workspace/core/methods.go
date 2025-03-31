package core

import (
	"context"
	"time"

	"github.com/auvitly/go-tools/standard/workspace/entity"
	"github.com/auvitly/go-tools/standard/workspace/storage"
	"github.com/auvitly/go-tools/stderrs"
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
		stderrD := c.dependencies.SessionStorage.Drop(ctx, storage.SessionDropParams{
			SessionID: session.ID,
		})
		if stderrD != nil {
			return nil, stderrD
		}

		return nil, stderr
	case stderr != nil:
		return nil, stderr
	default:
		return task, nil
	}
}

func (c *Core[T, M, S]) ReportState(ctx context.Context, params ReportStateParams[S]) *stderrs.Error {
	var ts = time.Now()

	task, stderr := c.dependencies.TaskStorage.Get(ctx, storage.TaskGetParams{
		TaskID: params.TaskID,
	})
	if stderr != nil {
		return stderr
	}

	switch {
	case task.SessionID == nil || task.AssignTS == nil:
		return stderrs.FailedPrecondition.SetMessage("task unassigned")
	case *task.SessionID != params.SessionID:
		return stderrs.InvalidArgument.SetMessage("worker session not match")
	case task.AssignTS.Sub(task.UpdatedTS) < c.config.TaskDowntime:
		stderr = c.dependencies.SessionStorage.Done(ctx, storage.SessionDoneParams{
			SessionID: *task.SessionID,
		})
		if stderr != nil {
			return stderr
		}

		_, stderr = c.dependencies.TaskStorage.Update(ctx, storage.TaskUpdateParams[S]{
			TaskID:       task.ID,
			StateData:    task.StateData,
			StatusCode:   task.StatusCode,
			Result:       task.Result,
			CatchLaterTS: nil,
			DoneTS:       nil,
			SessionID:    nil,
			AssignTS:     nil,
		})
		if stderr != nil {
			return stderr
		}

		return stderrs.FailedPrecondition.SetMessage("task unassigned")
	}

	var updateParams = storage.TaskUpdateParams[S]{
		TaskID:    params.TaskID,
		UpdatedTS: ts,
	}

	switch state := params.ReportState.(type) {
	case SetStateDone[S]:
		updateParams.AssignTS = nil
		updateParams.SessionID = nil
		updateParams.CatchLaterTS = nil
		updateParams.DoneTS = &ts
		updateParams.Result = &state.Result
		updateParams.StatusCode = &state.StatusCode
	case SetStateInWork:
		updateParams.StateData = state.StateData
	case SetStatePutOff:
		updateParams.StateData = state.StateData
		updateParams.CatchLaterTS = &state.CatchLaterAT
	default:
		return stderrs.InvalidArgument.SetMessage("not found state in report")
	}

	task, stderr = c.dependencies.TaskStorage.Update(ctx, updateParams)
	if stderr != nil {
		return stderr
	}

	if task.DoneTS == nil && task.CatchLaterTS == nil {
		return nil
	}

	stderr = c.dependencies.SessionStorage.Done(ctx, storage.SessionDoneParams{
		SessionID: params.SessionID,
	})
	if stderr != nil {
		return stderr
	}

	return nil
}

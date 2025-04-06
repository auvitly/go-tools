package service

import (
	"context"
	"time"

	"github.com/auvitly/go-tools/collection/stderrs"
	"github.com/auvitly/go-tools/standard/sdk/scheduler/entity/session"
	"github.com/auvitly/go-tools/standard/sdk/scheduler/entity/task"
	"github.com/auvitly/go-tools/standard/sdk/scheduler/entity/worker"
	"github.com/auvitly/go-tools/standard/utils/reflector"
	"github.com/google/uuid"
)

func (s *Service[T, W, S]) ReceiveTask(ctx context.Context, params ReceiveTaskParams[W]) (T, *stderrs.Error) {
	worker, stderr := s.dependencies.WorkerStorage.Save(ctx, params.Worker)
	if stderr != nil {
		return reflector.Nil[T](), stderr
	}

	list, stderr := s.dependencies.TaskStorage.List(ctx,
		task.ListFilterTypes{Types: params.Worker.Impl().Types},
		task.ListFilterLabels{WorkerLabels: params.Labels},
		task.ListFilterModes{Modes: params.Modes},
		task.ListFilterStatuses{Statuses: []task.Status{
			task.StatusCreated, task.StatusInProgress, task.StatusError, task.StatusCompensating,
		}},
		task.ListFilterWithoutAssigned{},
		task.ListFilterPagination{Limit: 1},
	)
	if stderr != nil {
		return reflector.Nil[T](), stderr
	}

	if len(list) == 0 {

	}

	workerSession, stderr := s.dependencies.SessionStorage.Receive(ctx, worker.Impl().ID)
	if stderr != nil {
		return reflector.Nil[T](), stderr
	}

	var (
		item = list[0]
		ts   = time.Now()
	)

	item.Impl().WorkerSessionID = &workerSession.Impl().ID
	item.Impl().WorkerAssignTS = &ts

	received, stderr := s.dependencies.TaskStorage.Push(ctx, item)
	if stderr != nil {
		stderr2 := s.dependencies.SessionStorage.Drop(ctx, workerSession.Impl().ID)
		if stderr2 != nil {
			return reflector.Nil[T](), stderr.EmbedErrors(stderr2)
		}

		return reflector.Nil[T](), stderr
	}

	return received, nil
}

func (s *Service[T, W, S]) ListTask(ctx context.Context, filters ...task.IsListFilter) ([]T, *stderrs.Error) {
	items, stderr := s.dependencies.TaskStorage.List(ctx, filters...)
	if stderr != nil {
		return nil, stderr
	}

	return items, nil
}

func (s *Service[T, W, S]) GetTask(ctx context.Context, id uuid.UUID) (result T, stderr *stderrs.Error) {
	founded, stderr := s.dependencies.TaskStorage.Get(ctx, id)
	if stderr != nil {
		return reflector.Nil[T](), stderr
	}

	return founded, nil
}

func (s *Service[T, W, S]) CancelTask(ctx context.Context, id uuid.UUID) (T, *stderrs.Error) {
	founded, stderr := s.dependencies.TaskStorage.Get(ctx, id)
	if stderr != nil {
		return reflector.Nil[T](), stderr
	}

	var impl = founded.Impl()

	if impl == nil {
		return reflector.Nil[T](), stderrs.Internal.SetMessage("not found implementation task by id")
	}

	switch impl.Status {
	case task.StatusCreated:
		founded.Impl().Status = task.StatusCanceled
	case task.StatusInProgress, task.StatusDone, task.StatusError:
		founded.Impl().Status = task.StatusCompensating
	case task.StatusCompensating, task.StatusCanceled:
		return founded, nil
	case task.StatusCompleted:
		return reflector.Nil[T](), stderrs.InvalidArgument.SetMessage("completed task cannot be cancelled")
	default:
		return reflector.Nil[T](), stderrs.Internal.SetMessage("unexpected task status '%s'", impl.Impl().Status)
	}

	pushed, stderr := s.dependencies.TaskStorage.Push(ctx, founded)
	if stderr != nil {
		return reflector.Nil[T](), stderr
	}

	return pushed, nil
}

func (s *Service[T, W, S]) CommitTask(ctx context.Context, id uuid.UUID) (T, *stderrs.Error) {
	founded, stderr := s.dependencies.TaskStorage.Get(ctx, id)
	if stderr != nil {
		return reflector.Nil[T](), stderr
	}

	if founded.Impl().Status != task.StatusDone {
		return reflector.Nil[T](), stderrs.InvalidArgument.
			SetMessage("task is not done, actual status '%s'", founded.Impl().Status)
	}

	founded.Impl().Status = task.StatusCompleted

	pushed, stderr := s.dependencies.TaskStorage.Push(ctx, founded)
	if stderr != nil {
		return reflector.Nil[T](), stderr
	}

	return pushed, nil
}

func (s *Service[T, W, S]) ReportTask(ctx context.Context, req ReportTaskRequest[T, W, S]) (*ReportTaskResponse[T], *stderrs.Error) {
	stderr := s.validateReportTaskRequest(req)
	if stderr != nil {
		return nil, stderr
	}

	founded, stderr := s.dependencies.TaskStorage.Get(ctx, req.TaskID)
	if stderr != nil {
		return nil, stderr
	}

	stderr = s.validateReportTaskAccess(req, founded)
	if stderr != nil {
		return nil, stderr
	}

	response, stderr := req.Event.makeReportTaskResponse(founded)
	if stderr != nil {
		return nil, stderr
	}

	if response.Task.Impl().WorkerSessionID == nil {
		stderr := s.dependencies.SessionStorage.Drop(ctx, req.Session.Impl().ID)
		if stderr != nil {
			return nil, stderr
		}
	}

	_, stderr = s.dependencies.TaskStorage.Push(ctx, response.Task)
	if stderr != nil {
		return nil, stderr
	}

	return response, nil
}

func (s *Service[T, W, S]) validateReportTaskRequest(req ReportTaskRequest[T, W, S]) *stderrs.Error {
	switch {
	case req.TaskID == uuid.Nil:
		return stderrs.Unauthenticated.SetMessage("not found task_id in report")
	case worker.IsWorker(req.Worker) == nil || req.Worker.Impl() == nil:
		return stderrs.Unauthenticated.SetMessage("not found worker in report")
	case session.IsSession(req.Session) == nil || req.Session.Impl() == nil:
		return stderrs.Unauthenticated.SetMessage("not found session in report")
	case req.Event == nil:
		return stderrs.InvalidArgument.SetMessage("not found event in report")
	default:
		return nil
	}
}

func (s *Service[T, W, S]) validateReportTaskAccess(req ReportTaskRequest[T, W, S], founded T) *stderrs.Error {
	switch {
	case task.IsTask(founded) == nil || founded.Impl() == nil:
		return stderrs.Internal.SetMessage("not found task implementation with id %s", req.TaskID.String())
	case founded.Impl().WorkerSessionID == nil:
		return stderrs.InvalidArgument.SetMessage("task %s not assigned", req.TaskID.String())
	case *founded.Impl().WorkerSessionID != req.Session.Impl().ID:
		return stderrs.PermissionDenied.SetMessage("task %s has been assigned to another performer", req.TaskID.String())
	default:
		return nil
	}
}

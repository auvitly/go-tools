package service

import (
	"time"

	"github.com/auvitly/go-tools/collection/stderrs"
	"github.com/auvitly/go-tools/standard/sdk/scheduler/entity/session"
	"github.com/auvitly/go-tools/standard/sdk/scheduler/entity/task"
	"github.com/auvitly/go-tools/standard/sdk/scheduler/entity/worker"
	"github.com/google/uuid"
)

type ReceiveTaskParams[W worker.IsWorker] struct {
	Worker W
	Modes  []task.Mode
	Labels map[string]string
}

type ReportTaskRequest[T task.IsTask, W worker.IsWorker, S session.IsSession] struct {
	TaskID  uuid.UUID
	Session S
	Worker  W
	Event   ReportTaskEvent[T]
}

type ReportTaskEvent[T task.IsTask] interface {
	makeReportTaskResponse(T) (*ReportTaskResponse[T], *stderrs.Error)
}

type ReportTaskEventInProcess struct {
	Data map[string]any
}

type ReportTaskEvenOnPending struct {
	Data      map[string]any
	PendingTS time.Time
}

type ReportTaskEventDone struct {
	Data   map[string]any
	Result map[string]any
}

type ReportTaskEventError struct {
	Data map[string]any
}

type ReportTaskEventCompensated struct {
	Data map[string]any
}

type ReportTaskResponse[T task.IsTask] struct {
	Task T
}

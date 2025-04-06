package scheduler

import (
	"context"
	"testing"

	"github.com/auvitly/go-tools/standard/sdk/scheduler/entity/session"
	"github.com/auvitly/go-tools/standard/sdk/scheduler/entity/task"
	"github.com/auvitly/go-tools/standard/sdk/scheduler/entity/worker"
	"github.com/auvitly/go-tools/standard/sdk/scheduler/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type MyTask struct {
	*task.Task
	Field string
}

type MySession struct {
	*session.Session
	Field string
}

type MyWorker struct {
	*worker.Worker
	Field string
}

func TestScheduler(t *testing.T) {
	var ctx context.Context

	s, stderr := service.New(service.Dependencies[*MyTask, *MyWorker, *MySession]{
		TaskStorage:    nil,
		WorkerStorage:  nil,
		SessionStorage: nil,
	})
	require.Nil(t, stderr)

	prepeared, stderr := task.New(task.NewParams{
		ID:        uuid.New(),
		ParentID:  nil,
		Type:      "TYPE_TEST",
		Mode:      task.ModeAsync,
		PendingTS: nil,
		Arguments: map[string]any{},
		Labels:    map[string]string{},
	})
	require.Nil(t, stderr)

	s.CreateTask(ctx, &MyTask{
		Task:  prepeared,
		Field: "field",
	})
}

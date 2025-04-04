package workspace_test

import (
	"context"
	"testing"
	"time"

	"github.com/auvitly/go-tools/standard/sdk/scheduler/service"
	"github.com/auvitly/go-tools/standard/sdk/scheduler/storage/inmemory"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

func TestCore(t *testing.T) {
	var ctx = context.Background()

	type (
		Type int
		Mode int
	)

	workspace, stderr := service.New(
		ctx,
		service.Dependencies[Type, Mode, codes.Code]{
			TaskStorage: inmemory.NewTaskStorage[Type, Mode, codes.Code](inmemory.TaskConfig{
				DeleteCompleted: true,
			}),
			WorkerStorage: inmemory.NewWorkerStorage[Type](),
			SessionStorage: inmemory.NewSessionStorage(inmemory.SessionConfig{
				DeleteCompleted: true,
			}),
		},
		service.Config{
			TaskDowntime:    time.Second,
			PullingInterval: time.Second,
		},
	)
	require.Nil(t, stderr)

	_, stderr = workspace.CreateTask(ctx, service.CreateTaskParams[Type, Mode]{
		Type: 1,
		Mode: 1,
		Args: nil,
		Labels: map[string]string{
			"worker": "A",
		},
	})
	require.Nil(t, stderr)

	task, stderr := workspace.ReceiveTask(ctx, service.ReceiveTaskParams[Type]{
		WorkerID: uuid.New(),
		Type:     1,
		Version:  "version",
		Labels: map[string]string{
			"worker": "A",
		},
	})
	require.Nil(t, stderr)

	time.Sleep(time.Second)

	stderr = workspace.ReportState(ctx, service.ReportStateParams[codes.Code]{
		TaskID:    task.ID,
		SessionID: *task.SessionID,
		ReportState: service.SetStatePutOff{
			StateData:    map[string]any{},
			CatchLaterAT: time.Now().Add(time.Second),
		},
	})
	require.NotNil(t, stderr)

	task, stderr = workspace.ReceiveTask(ctx, service.ReceiveTaskParams[Type]{
		WorkerID: uuid.New(),
		Type:     1,
		Version:  "version",
		Labels: map[string]string{
			"worker": "A",
		},
	})
	require.Nil(t, stderr)

	stderr = workspace.ReportState(ctx, service.ReportStateParams[codes.Code]{
		TaskID:    task.ID,
		SessionID: *task.SessionID,
		ReportState: service.SetStatePutOff{
			StateData:    map[string]any{},
			CatchLaterAT: time.Now().Add(time.Second),
		},
	})
	require.Nil(t, stderr)

	t.Log(task)
}

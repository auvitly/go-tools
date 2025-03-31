package workspace_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/auvitly/go-tools/standard/workspace/core"
	"github.com/auvitly/go-tools/standard/workspace/storage/inmemory"
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

	workspace, stderr := core.New(
		core.Dependencies[Type, Mode, codes.Code]{
			TaskStorage:    inmemory.NewTaskStorage[Type, Mode, codes.Code](inmemory.TaskConfig{DeleteCompleted: true}),
			WorkerStorage:  inmemory.NewWorkerStorage[Type](),
			SessionStorage: inmemory.NewSessionStorage(inmemory.SessionConfig{DeleteCompleted: true}),
		},
		core.Config{
			TaskDowntime:    time.Second,
			PullingInterval: time.Second,
		},
	)
	require.Nil(t, stderr)

	_, stderr = workspace.CreateTask(ctx, core.CreateTaskParams[Type, Mode]{
		Type:   1,
		Mode:   1,
		Status: "created",
		Args:   nil,
		Labels: map[string]string{
			"worker": "A",
		},
	})
	require.Nil(t, stderr)

	task, stderr := workspace.ReceiveTask(ctx, core.ReceiveTaskParams[Type]{
		WorkerID: uuid.New(),
		Type:     1,
		Version:  "version",
		Labels: map[string]string{
			"worker": "A",
		},
	})
	require.Nil(t, stderr)

	time.Sleep(time.Second)

	stderr = workspace.ReportState(ctx, core.ReportStateParams[codes.Code]{
		TaskID:    task.ID,
		SessionID: *task.SessionID,
		ReportState: core.SetStatePutOff{
			StateData:    json.RawMessage([]byte("{}")),
			CatchLaterAT: time.Now().Add(time.Second),
		},
	})
	require.NotNil(t, stderr)

	task, stderr = workspace.ReceiveTask(ctx, core.ReceiveTaskParams[Type]{
		WorkerID: uuid.New(),
		Type:     1,
		Version:  "version",
		Labels: map[string]string{
			"key": "value",
		},
	})
	require.Nil(t, stderr)

	stderr = workspace.ReportState(ctx, core.ReportStateParams[codes.Code]{
		TaskID:    task.ID,
		SessionID: *task.SessionID,
		ReportState: core.SetStatePutOff{
			StateData:    json.RawMessage([]byte("{}")),
			CatchLaterAT: time.Now().Add(time.Second),
		},
	})
	require.Nil(t, stderr)

	t.Log(task)
}

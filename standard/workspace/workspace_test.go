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
		Type       int
		Mode       int
		StatusCode int
	)

	workspace, stderr := core.New(
		core.Dependencies[Type, Mode, StatusCode]{
			TaskStorage:    inmemory.NewTaskStorage[Type, Mode, StatusCode](),
			WorkerStorage:  inmemory.NewWorkerStorage[Type](),
			SessionStorage: inmemory.NewSessionStorage(),
		},
		core.Config{
			SessionDecayTime: time.Hour,
			PullingInterval:  time.Second,
		},
	)
	require.Nil(t, stderr)

	_, stderr = workspace.CreateTask(ctx, core.CreateTaskParams[Type, Mode]{
		Type:   1,
		Mode:   1,
		Status: "created",
		Args:   nil,
		Labels: map[string]string{
			"key": "value",
		},
	})
	require.Nil(t, stderr)

	var workerID = uuid.New()

	task, stderr := workspace.ReceiveTask(ctx, core.ReceiveTaskParams[Type]{
		WorkerID: workerID,
		Type:     1,
		Version:  "version",
		Labels: map[string]string{
			"key": "value",
		},
	})
	require.Nil(t, stderr)

	var result = json.RawMessage([]byte("{}"))

	stderr = workspace.ReportState(ctx, core.ReportStateParams[StatusCode]{
		TaskID:    task.ID,
		SessionID: *task.SessionID,
		ReportState: core.SetStateDone[StatusCode]{
			StatusCode: StatusCode(codes.OK),
			Result:     result,
		},
	})
	require.Nil(t, stderr)

	t.Log(task)
}

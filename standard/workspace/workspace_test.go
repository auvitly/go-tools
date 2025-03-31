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
)

func TestCore(t *testing.T) {
	var ctx = context.Background()

	workspace, stderr := core.New(
		core.Dependencies[int, int]{
			TaskStorage:    inmemory.NewTaskStorage[int, int](),
			WorkerStorage:  inmemory.NewWorkerStorage[int](),
			SessionStorage: inmemory.NewSessionStorage[int](),
		},
		core.Config{
			SessionDecayTime: time.Hour,
			PullingInterval:  time.Second,
		},
	)
	require.Nil(t, stderr)

	_, stderr = workspace.CreateTask(ctx, core.CreateTaskParams[int, int]{
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

	task, stderr := workspace.ReceiveTask(ctx, core.ReceiveTaskParams[int]{
		WorkerID: workerID,
		Type:     1,
		Version:  "version",
		Labels: map[string]string{
			"key": "value",
		},
	})
	require.Nil(t, stderr)

	var result = json.RawMessage([]byte("{}"))

	stderr = workspace.SetState(ctx, core.SetStateParams[int]{
		TaskID:    task.ID,
		SessionID: *task.SessionID,
		Status:    "done",
		Result:    &result,
	})
	require.Nil(t, stderr)

	t.Log(task)
}

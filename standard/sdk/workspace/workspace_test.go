package workspace_test

import (
	"testing"

	"github.com/auvitly/go-tools/standard/sdk/workspace"
	"github.com/stretchr/testify/require"
)

type Stage int

const (
	StageInit Stage = iota
	StageError
	StageDone
)

func TestWorkspace(t *testing.T) {
	type Struct struct {
		A string `json:"a"`
	}

	var ws = workspace.New(StageInit, "my init stage")

	workspace.Store(ws,
		workspace.KV{Key: "string", Value: "string"},
		workspace.KV{Key: "int", Value: 0},
		workspace.KV{Key: "struct", Value: Struct{A: "a"}},
	)

	workspace.SetStage(ws, StageError, "it's error!")

	ws2, stderr := workspace.FromMap[Stage](workspace.ToMap(ws))
	require.Nil(t, stderr)

	str, stderr := workspace.Load[Struct](ws2, "struct")
	require.Nil(t, stderr)

	t.Log(str)
}

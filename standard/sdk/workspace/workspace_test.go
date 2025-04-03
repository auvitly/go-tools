package workspace_test

import (
	"testing"

	"github.com/auvitly/go-tools/standard/sdk/workspace"
	"github.com/stretchr/testify/require"
)

func TestWorkspace(t *testing.T) {
	type Stage int
	type Struct struct {
		A string `json:"a"`
	}

	var ws = workspace.New[Stage](0, "init")

	workspace.Store(ws, "string", "string")
	workspace.Store(ws, "int", 0)
	workspace.Store(ws, "struct", Struct{A: "a"})

	ws2, stderr := workspace.FromMap[Stage](workspace.ToMap(ws))
	require.Nil(t, stderr)

	str, stderr := workspace.Load[Struct](ws2, "struct")
	require.Nil(t, stderr)

	t.Log(str)
}

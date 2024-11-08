package labfunc

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

// JSONMarshal - json.Marshal without error.
func JSONMarshal(t *testing.T, v any) []byte {
	t.Helper()

	data, err := json.Marshal(v)
	require.NoError(t, err)

	return data
}

// JSONUnmarshal - json.Unmarshal without error.
func JSONUnmarshal[O any](t *testing.T, data []byte) (object O) {
	t.Helper()

	err := json.Unmarshal(data, &object)
	require.NoError(t, err)

	return object
}

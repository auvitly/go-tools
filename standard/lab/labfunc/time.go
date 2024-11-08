package labfunc

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// ParseTime - time.ParseTime without error. Returns pointer.
func ParseTime(t *testing.T, layout, s string) *time.Time {
	t.Helper()

	ts, err := time.Parse(layout, s)
	require.NoError(t, err)

	return &ts
}

// ParseDuration - time.ParseDuration without error. Returns pointer.
func ParseDuration(t *testing.T, s string) *time.Duration {
	t.Helper()

	dur, err := time.ParseDuration(s)
	require.NoError(t, err)

	return &dur
}

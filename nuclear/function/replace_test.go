package function_test

import (
	"github.com/auvitly/go-tools/nuclear/function"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestReplace(t *testing.T) {
	t.Parallel()

	var ts = time.Now()

	var (
		oldTimeFunc func() time.Time
		newTimeFunc = func() time.Time {
			return ts
		}
	)

	patch := function.Replace(time.Now, newTimeFunc, &oldTimeFunc)
	require.NotNil(t, patch)

	assert.Equal(t, ts, time.Now())

	patch.Unpatch()
	assert.NotEqual(t, ts, time.Now())
}

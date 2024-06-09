package function_test

import (
	"github.com/auvitly/go-tools/nuclear/function"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestReplace(t *testing.T) {
	var ts = time.Now()

	var (
		oldTimeFunc func() time.Time
		newTimeFunc = func() time.Time {
			return ts
		}
	)

	patch := function.Replace(time.Now, newTimeFunc, &oldTimeFunc)
	require.NotNil(t, patch)

	assert.Equal(t, time.Now(), ts, "equal")

	patch.Unpatch()
	time.Sleep(time.Second)
	assert.NotEqual(t, time.Now(), ts, "not equal")
}

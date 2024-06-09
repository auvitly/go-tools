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

	var fn = func() time.Time {
		return ts
	}

	time.Sleep(time.Second)

	for i := 0; i < 100; i++ {
		patch := function.Replace(time.Now, fn)
		require.NotNil(t, patch)

		assert.Equal(t, time.Now(), ts, "equal")

		patch.Unpatch()
		assert.NotEqual(t, time.Now(), ts, "not equal")
	}
}

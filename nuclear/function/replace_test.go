package function_test

import (
	"github.com/auvitly/go-tools/nuclear/function"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestReplace(t *testing.T) {
	t.Parallel()

	var (
		oldTimeFunc func() time.Time
		newTimeFunc = func() time.Time {
			return time.Date(2000, 0, 0, 0, 0, 0, 0, time.Local)
		}
	)

	oldTimeFunc = function.Replace(time.Now, newTimeFunc).OldImpl()

	a := oldTimeFunc()
	b := time.Now()

	assert.NotEqual(t, a, b)
}

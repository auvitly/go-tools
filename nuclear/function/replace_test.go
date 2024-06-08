package function_test

import (
	"github.com/auvitly/go-tools/nuclear/function"
	"testing"
	"time"
)

func TestReplace(t *testing.T) {
	t.Parallel()

	var (
		old  func() time.Time
		impl = func() time.Time {
			return time.Unix(0, 0)
		}
	)

	function.Replace(time.Now, impl, &old)

	t.Log(time.Now())
}

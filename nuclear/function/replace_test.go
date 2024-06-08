package function_test

import (
	"github.com/auvitly/go-tools/nuclear/function"
	"testing"
)

func Sum(a, b int) int {
	return a + b
}

func TestReplace(t *testing.T) {
	t.Parallel()

	var (
		old  func(a, b int) int
		impl = func(a, b int) int {
			return 1
		}
	)

	function.Replace(Sum, impl, &old)

	t.Log(Sum(0, 0))
}

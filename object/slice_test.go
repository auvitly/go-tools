package object_test

import (
	"github.com/auvitly/go-tools/object"
	"testing"
)

func TestSlice(t *testing.T) {
	var m = object.Map[int, int]{3: 1}.
		ScanKV([]int{1, 2}, []int{3, 4})

	t.Log(m)
}

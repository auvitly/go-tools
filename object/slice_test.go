package object_test

import (
	"github.com/auvitly/go-tools/object"
	"testing"
)

func TestSlice(t *testing.T) {
	var listA = object.SliceOf[int]{1, 2, 3, 4}
	var listC = make(object.SliceOf[int], 4)

	listD := listA.Join(listA, listA, listA)

	listB := listA.Clone().Reverse().Copy(listC).ForEach(func(_ int, item *int) { *item = *item * 2 })

	t.Log(listA, listB, listC, listD)
}

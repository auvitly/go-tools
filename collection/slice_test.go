package collection_test

import (
	"github.com/auvitly/go-tools/collection/action"
	"github.com/auvitly/go-tools/collection/object"
	"github.com/auvitly/go-tools/collection/order"
	"testing"
)

func TestSlice(t *testing.T) {
	var items = object.Slice[int]{1, 2, 2, 0, 1, 3, 0, 3}

	values := items.Clone().
		Action(action.FilterUnique[int]).
		Sort(order.ASC[int])

	t.Log(items, values)
}

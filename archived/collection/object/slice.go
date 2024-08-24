package object

import (
	"github.com/auvitly/go-tools/collection/order"
	"slices"
)

// Slice - slice model of type E.
type (
	Slice[E any]          []E
	SliceAction[I any]    func(slice Slice[I]) Slice[I]
	SliceCondition[I any] func(index int, item I) bool
)

// SliceOf - method for converting a model into a slice.
func SliceOf[T any](s []T) Slice[T] { return s }

func (s Slice[E]) Len() int { return len(s) }

func (s Slice[E]) Cap() int { return cap(s) }

func (s Slice[E]) Append(elems ...E) Slice[E] { return append(s, elems...) }

func (s Slice[E]) AppendSlice(elems ...Slice[E]) Slice[E] {
	var capacity = s.Cap()

	for _, elem := range elems {
		capacity += elem.Cap()
	}

	var result = make(Slice[E], 0, capacity)

	result = result.Append(s...)

	for _, elem := range elems {
		result = result.Append(elem...)
	}

	return result
}

func (s Slice[E]) CopyFrom(src Slice[E]) Slice[E] { copy(s, src); return s }

func (s Slice[E]) CopyTo(dst Slice[E]) Slice[E] { copy(dst, s); return s }

func (s Slice[E]) Sort(rule order.Sort[E]) Slice[E] { slices.SortFunc(s, rule); return s }

func (s Slice[E]) Reverse() Slice[E] { slices.Reverse(s); return s }

func (s Slice[E]) Clip() Slice[E] { return s[:len(s):len(s)] }

func (s Slice[E]) Clone() Slice[E] { return make(Slice[E], s.Len(), s.Cap()).CopyFrom(s) }

func (s Slice[E]) Chunk(size uint) Chunk[E] {
	if size == 0 {
		size = 1
	}

	var num = s.Len() / int(size)

	if num == 0 {
		return Chunk[E]{s}
	}

	if len(s)%int(size) != 0 {
		num++
	}

	var result = make(Chunk[E], 0, num)

	for i := 0; i < num; i++ {
		var last = (i + 1) * int(size)

		if last > s.Len() {
			last = s.Len()
		}

		result = append(result, s[i*int(size):last])
	}

	return result
}

func (s Slice[E]) Delete(i, j int) Slice[E] {
	_ = s[i:j]

	return append(s[:i], s[j:]...)
}

func (s Slice[E]) DeleteFunc(del func(E) bool) Slice[E] {
	return slices.DeleteFunc(s, del)
}

func (s Slice[E]) Is(conditions ...SliceCondition[E]) (ok bool) {
	if len(conditions) == 0 {
		return false
	}

	var fn SliceCondition[E] = func(index int, item E) bool {
		for i := range conditions {
			if !conditions[i](index, item) {
				return false
			}
		}

		return true
	}

	for i, item := range s {
		if !fn(i, item) {
			return false
		}
	}

	return true
}

func (s Slice[E]) Action(actions ...SliceAction[E]) Slice[E] {
	var result = s

	for i := range actions {
		result = actions[i](result)
	}

	return result
}

func (s Slice[E]) ForEach(fn func(index int, item *E)) Slice[E] {
	if fn == nil {
		return s
	}

	for index := range s {
		fn(index, &s[index])
	}

	return s
}

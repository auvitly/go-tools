package object

import "slices"

// Slice - interface for explicitly passing slices.
type Slice interface {
	implSlice()
}

// SliceOf - slice model for type T.
type SliceOf[T any] []T

func (s SliceOf[T]) implSlice()                            {}
func (s SliceOf[T]) Append(elems ...T) SliceOf[T]          { return append(s, elems...) }
func (s SliceOf[T]) Len() int                              { return len(s) }
func (s SliceOf[T]) Cap() int                              { return cap(s) }
func (s SliceOf[T]) Clone() SliceOf[T]                     { return slices.Clone(s) }
func (s SliceOf[T]) Scan(src SliceOf[T]) SliceOf[T]        { copy(s, src); return s }
func (s SliceOf[T]) Copy(dst SliceOf[T]) SliceOf[T]        { copy(dst, s); return s }
func (s SliceOf[T]) Sort(rule func(a, b T) int) SliceOf[T] { slices.SortFunc(s, rule); return s }
func (s SliceOf[T]) Reverse() SliceOf[T]                   { slices.Reverse(s); return s }
func (s SliceOf[T]) Clip() SliceOf[T]                      { slices.Clip(s); return s }

func (s SliceOf[T]) Delete(condition func(item T) bool) SliceOf[T] {
	if condition == nil {
		return s
	}

	slices.DeleteFunc(s, condition)

	return s
}

func (s SliceOf[T]) Join(elems ...SliceOf[T]) SliceOf[T] {
	var capacity = s.Cap()

	for _, elem := range elems {
		capacity += elem.Cap()
	}

	var result = make(SliceOf[T], 0, capacity)

	result = result.Append(s...)

	for _, elem := range elems {
		result = result.Append(elem...)
	}

	return result
}

func (s SliceOf[T]) IsCondition(condition func(index int, item T) bool) (ok bool) {
	if condition == nil {
		return false
	}

	for index, item := range s {
		if condition(index, item) {
			return true
		}
	}

	return false
}

func (s SliceOf[T]) Index(condition func(index int, item T) bool) (index int) {
	if condition == nil {
		return -1
	}

	for index = range s {
		if condition(index, s[index]) {
			return index
		}
	}

	return -1
}

func (s SliceOf[T]) Find(condition func(index int, item T) bool) (index int, item T) {
	if condition == nil {
		return -1, item
	}

	for index = range s {
		if condition(index, s[index]) {
			return index, s[index]
		}
	}

	return -1, item
}

func (s SliceOf[T]) ForEach(fn func(index int, item *T)) SliceOf[T] {
	if fn == nil {
		return s
	}

	for index := range s {
		fn(index, &s[index])
	}

	return s
}

package object

import "slices"

// Slice - slice model of type T.
type Slice[T any] []T

func (s Slice[T]) Len() int { return len(s) }

func (s Slice[T]) Cap() int { return cap(s) }

func (s Slice[T]) Append(elems ...T) Slice[T] { return append(s, elems...) }

func (s Slice[T]) AppendSlice(elems ...Slice[T]) Slice[T] {
	var capacity = s.Cap()

	for _, elem := range elems {
		capacity += elem.Cap()
	}

	var result = make(Slice[T], 0, capacity)

	result = result.Append(s...)

	for _, elem := range elems {
		result = result.Append(elem...)
	}

	return result
}

func (s Slice[T]) CopyFrom(src Slice[T]) Slice[T] { copy(s, src); return s }

func (s Slice[T]) CopyTo(dst Slice[T]) Slice[T] { copy(dst, s); return s }

func (s Slice[T]) Sort(rule func(a, b T) int) Slice[T] { slices.SortFunc(s, rule); return s }

func (s Slice[T]) Reverse() Slice[T] { slices.Reverse(s); return s }

func (s Slice[T]) Clip() Slice[T] { return s[:len(s):len(s)] }

func (s Slice[T]) Clone() Slice[T] { return make(Slice[T], s.Len(), s.Cap()).CopyFrom(s) }

func (s Slice[T]) Chunk(size uint) Chunk[T] {
	if size == 0 {
		size = 1
	}

	var num = s.Len() / int(size)

	if num == 0 {
		return Chunk[T]{s}
	}

	if len(s)%int(size) != 0 {
		num++
	}

	var result = make(Chunk[T], 0, num)

	for i := 0; i < num; i++ {
		var last = (i + 1) * int(size)

		if last > s.Len() {
			last = s.Len()
		}

		result = append(result, s[i*int(size):last])
	}

	return result
}

func (s Slice[T]) Delete(condition func(item T) bool) Slice[T] {
	if condition == nil {
		return s
	}

	slices.DeleteFunc(s, condition)

	return s
}

func (s Slice[T]) IsCondition(condition func(index int, item T) bool) (ok bool) {
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

func (s Slice[T]) Index(condition func(index int, item T) bool) (index int) {
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

func (s Slice[T]) Find(condition func(index int, item T) bool) (index int, item T) {
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

func (s Slice[T]) ForEach(fn func(index int, item *T)) Slice[T] {
	if fn == nil {
		return s
	}

	for index := range s {
		fn(index, &s[index])
	}

	return s
}

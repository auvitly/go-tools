package object

import "slices"

// Chunk - slice of elements split into groups the length of size.
type Chunk[T any] []Slice[T]

func (c Chunk[T]) Append(elems ...Slice[T]) Chunk[T]          { return append(c, elems...) }
func (c Chunk[T]) Sort(rule func(a, b Slice[T]) int) Chunk[T] { slices.SortFunc(c, rule); return c }
func (c Chunk[T]) Reverse() Chunk[T]                          { slices.Reverse(c); return c }

func (c Chunk[T]) SortItems(rule func(a, b T) int) Chunk[T] {
	return c.ForEach(func(index int, item *Slice[T]) {
		item.Sort(rule)
	})
}

func (c Chunk[T]) ReverseItems() Chunk[T] {
	return c.ForEach(func(index int, item *Slice[T]) {
		item.Reverse()
	})
}

func (c Chunk[T]) Clone() Chunk[T] {
	var result = make(Chunk[T], 0, cap(c))

	c.ForEach(func(_ int, item *Slice[T]) {
		result = result.Append(item.Clone())
	})

	return result
}

func (c Chunk[T]) Len() int {
	var length int

	for i := range c {
		length += len(c[i])
	}

	return length
}

func (c Chunk[T]) ForEach(fn func(index int, item *Slice[T])) Chunk[T] {
	if fn == nil {
		return c
	}

	for index := range c {
		fn(index, &c[index])
	}

	return c
}

func (c Chunk[T]) Redistribute(size uint) Chunk[T] {
	var items Slice[T]

	for i := 0; i < len(c); i++ {
		items = items.Append(c[i]...)
	}

	return items.Chunk(size)
}

func (c Chunk[T]) Join() Slice[T] {
	var capacity int

	for _, elem := range c {
		capacity += cap(elem)
	}

	var result = make(Slice[T], 0, capacity)

	for _, elem := range c {
		result = result.Append(elem...)
	}

	return result
}

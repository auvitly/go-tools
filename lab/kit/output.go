package kit

type Out[E any] struct {
	Expected E
	Error    error
}

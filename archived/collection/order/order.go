package order

import (
	"cmp"
)

type Sort[T any] func(a, b T) int

func ASC[T cmp.Ordered](a, b T) int {
	return cmp.Compare(a, b)
}

func DESC[T cmp.Ordered](a, b T) int {
	return ASC(a, b) * -1
}

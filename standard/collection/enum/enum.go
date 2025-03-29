package enum

import (
	"slices"

	"github.com/auvitly/go-tools/stderrs"
)

// Enum - .
type Enum[T comparable] []T

// Contains - .
func (e *Enum[T]) Contains(value T) bool {
	return slices.Contains(*e, value)
}

// Registry - .
func (e *Enum[T]) Registry(value T) (T, *stderrs.Error) {
	if e.Contains(value) {
		return *new(T), stderrs.InvalidArgument.
			SetMessage("value %v already exists in enum %T", value, value)
	}

	*e = append(*e, value)

	return value, nil
}

// MustRegistry - .
func (e *Enum[T]) MustRegistry(value T) T {
	result, err := e.Registry(value)
	if err != nil {
		panic(err)
	}

	return result
}

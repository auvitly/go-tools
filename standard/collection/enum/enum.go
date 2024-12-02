package enum

import (
	"errors"
)

// Enum - .
type Enum[T comparable] []T

// ErrAlreadyExists - .
var ErrAlreadyExists = errors.New("already exists")

// Contains - .
func (e *Enum[T]) Contains(value T) bool {
	for _, item := range *e {
		if item == value {
			return true
		}
	}

	return false
}

// Registry - .
func (e *Enum[T]) Registry(value T) (T, error) {
	if e.Contains(value) {
		return *new(T), ErrAlreadyExists
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

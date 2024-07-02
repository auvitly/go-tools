package object

import (
	"errors"
	"reflect"
)

// Enum - enumeration model for type T.
type Enum[T any] []T

// ErrAlreadyExists - the object is already contained in the enumeration.
var ErrAlreadyExists = errors.New("already exists")

// Add - add an object to the enumeration.
func (e *Enum[T]) Add(value T) (T, error) {
	for _, item := range *e {
		if reflect.DeepEqual(item, value) {
			return *new(T), ErrAlreadyExists
		}
	}

	*e = append(*e, value)

	return value, nil
}

// Contains - enum contains element.
func (e *Enum[T]) Contains(value T) bool {
	if len(*e) == 0 {
		return false
	}

	for _, item := range *e {
		if reflect.DeepEqual(item, value) {
			return true
		}
	}

	return false
}

// MustAdd - add an object to the enumeration. Panics if already exists.
func (e *Enum[T]) MustAdd(value T) T {
	item, err := e.Add(value)
	if err != nil {
		panic(err)
	}

	return item
}

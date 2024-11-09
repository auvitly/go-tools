package lab

import (
	"reflect"
)

// Pointer - returns pointer on copy value.
func Pointer[T any](v T) *T {
	return &v
}

// Value -returns the first value if the error is nil.
// Signature: func(*) (V, error).
func Value[V any, E error](value1 V, err E) V {
	if reflect.ValueOf(err).IsValid() && !reflect.ValueOf(err).IsNil() {
		panic(err)
	}

	return value1
}

// Error - returns the error.
// Signature: func(*) (V, error).
func Error[V any, E error](_ V, err E) E {
	return err
}

// FirstValue - returns the first result if the error is nil.
// Signature: func(*) (V1, V2, error).
func FirstValue[V1, V2 any, E error](value1 V1, _ V2, err E) V1 {
	if reflect.ValueOf(err).IsValid() && !reflect.ValueOf(err).IsNil() {
		panic(err)
	}

	return value1
}

// SecondValue - returns the second result if the error is nil.
// Signature: func(*) (V1, V2, error).
func SecondValue[V1, V2 any, E error](_ V1, value2 V2, err E) V2 {
	if reflect.ValueOf(err).IsValid() && !reflect.ValueOf(err).IsNil() {
		panic(err)
	}

	return value2
}

// ThirdError - returns the error from signature func(*) (V1, V2, error).
// Signature: func(*) (V1, V2, error).
func ThirdError[V1, V2 any, E error](_ V1, _ V2, err E) E {
	return err
}

package lab

import "fmt"

// Pointer - returns pointer on copy value.
func Pointer[T any](v T) *T {
	return &v
}

type ReturnStatement int

const (
	FirstValue ReturnStatement = iota
	SecondValue
	ThirdValue
	FourthValue
)

// Return - returns function with return value if the error is nil.
func Return[T any](args ...any) func(i ReturnStatement) T {
	return func(i ReturnStatement) T {
		switch {
		case len(args) == 0:
			panic("not found return values")
		case len(args) > int(i):
			if t, ok := args[i].(T); !ok {
				panic(fmt.Sprintf("value by statement %d is %T not %T", i, args[i], t))
			} else {
				return t
			}
		default:
			panic("out of range from args")
		}
	}
}

// First - returns the first value.
// Signature: func(*) (V1, V2).
func First[V1 any, V2 any](value V1, _ V2) V1 {
	return value
}

// Second - returns the second value.
// Signature: func(*) (V1, V2).
func Second[V1 any, V2 any](_ V1, value V2) V2 {
	return value
}

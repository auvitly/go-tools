package lab

// Pointer - returns pointer on copy value.
func Pointer[T any](v T) *T {
	return &v
}

type ReturnStatement int

const (
	FirstValue  ReturnStatement = 0
	SecondValue ReturnStatement = 1
	ThirdValue  ReturnStatement = 2
)

// Return - returns function with return value if the error is nil.
// Signature: func(*) (V1, V2).
func Return[T any](args ...any) func(i ReturnStatement) T {
	return func(i ReturnStatement) (t T) {
		switch {
		case len(args) == 0:
			panic("not found return values")
		case len(args) > int(i):
			return args[i].(T)
		default:
			panic("out of range from args")
		}
	}
}

// First - returns the first value if the error is nil.
// Signature: func(*) (V1, V2).
func First[V1 any, V2 any](value V1, _ V2) V1 {
	return value
}

// Second - returns the error.
// Signature: func(*) (V1, V2).
func Second[V1 any, V2 any](_ V1, value V2) V2 {
	return value
}

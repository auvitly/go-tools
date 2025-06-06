package lab

// Pointer - returns pointer on copy value.
func Pointer[T any](v T) *T {
	return &v
}

// PullOut - convert any values to slice any.
func PullOut[T any](index int) func(args ...any) T {
	return func(args ...any) T {
		return args[index].(T)
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

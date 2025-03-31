package lab

// Pointer - returns pointer on copy value.
func Pointer[T any](v T) *T {
	return &v
}

// Get - convert any values to slice any.
func Get[T any](args ...any) func(index int) T {
	return func(index int) T {
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

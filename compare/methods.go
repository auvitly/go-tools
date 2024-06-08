package compare

// Equal - equality check.
func Equal[T ComparableEqual](a, b T) bool {
	return a.Equal(b)
}

// NotEqual - inequality check.
func NotEqual[T ComparableEqual](a, b T) bool {
	return !a.Equal(b)
}

// Greater - checking whether value a is greater than value b.
func Greater[T ComparableGreater](a, b T) bool {
	return a.GreaterThan(b)
}

// NotGreater - checking value a is not greater than value b.
func NotGreater[T ComparableGreater](a, b T) bool {
	return !a.GreaterThan(b)
}

// GreaterOrEqual - checking value a is greater than or equal to value b.
func GreaterOrEqual[T interface {
	ComparableGreater
	ComparableEqual
}](a, b T) bool {
	return a.Equal(b) || a.GreaterThan(b)
}

// Less - checking whether value a is less than value b.
func Less[T ComparableLess](a, b T) bool {
	return a.LessThan(b)
}

// NotLess - checking whether value a is not less than value b.
func NotLess[T ComparableLess](a, b T) bool {
	return !a.LessThan(b)
}

func LessOrEqual[T interface {
	ComparableLess
	ComparableEqual
}](a, b T) bool {
	return a.Equal(b) || a.LessThan(b)
}

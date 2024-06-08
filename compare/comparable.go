package compare

// ComparableEqual - comparison type.
type ComparableEqual interface {
	Equal(c ComparableEqual) bool
}

// ComparableLess - comparison type.
type ComparableLess interface {
	LessThan(c ComparableLess) bool
}

// ComparableGreater - comparison type.
type ComparableGreater interface {
	GreaterThan(c ComparableGreater) bool
}

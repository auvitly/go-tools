package where

import (
	"cmp"
	"github.com/auvitly/go-tools/collection/object"
	"slices"
)

func IndexGreat[I any](i int) object.SliceCondition[I] {
	return func(index int, _ I) bool {
		if index > i {
			return true
		}

		return false
	}
}

func IndexLess[I any](i int) object.SliceCondition[I] {
	return func(index int, _ I) bool {
		if index < i {
			return true
		}

		return false
	}
}

func IndexRange[I any](i, j int) object.SliceCondition[I] {
	return func(index int, _ I) bool {
		if index >= i && index < j {
			return true
		}

		return false
	}
}

func ElemGreat[I cmp.Ordered](value I) object.SliceCondition[I] {
	return func(index int, item I) bool {
		if item > value {
			return true
		}

		return false
	}
}

func ElemGreatOrEqual[I cmp.Ordered](value I) object.SliceCondition[I] {
	return func(index int, item I) bool {
		if item >= value {
			return true
		}

		return false
	}
}

func ElemLess[I cmp.Ordered](value I) object.SliceCondition[I] {
	return func(index int, item I) bool {
		if item > value {
			return true
		}

		return false
	}
}

func ElemLessOrEqual[I cmp.Ordered](value I) object.SliceCondition[I] {
	return func(index int, item I) bool {
		if item >= value {
			return true
		}

		return false
	}
}

func ElemRange[I cmp.Ordered](i, j I) object.SliceCondition[I] {
	return func(index int, item I) bool {
		if item >= i && item < j {
			return true
		}

		return false
	}
}

func ElemEqual[I comparable](items ...I) object.SliceCondition[I] {
	return func(_ int, item I) bool {
		if slices.Contains(items, item) {
			return true
		}

		return false
	}
}

func ElemNotEqual[I comparable](items ...I) object.SliceCondition[I] {
	return func(_ int, item I) bool {
		if !slices.Contains(items, item) {
			return true
		}

		return false
	}
}

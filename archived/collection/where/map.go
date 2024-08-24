package where

import (
	"cmp"
	"github.com/auvitly/go-tools/collection/object"
	"slices"
)

func KeyEqual[K comparable, V any](keys ...K) object.MapCondition[K, V] {
	return func(key K, value V) bool {
		if slices.Contains(keys, key) {
			return true
		}

		return false
	}
}

func KeyNotEqual[K comparable, V any](keys ...K) object.MapCondition[K, V] {
	return func(key K, value V) bool {
		if !slices.Contains(keys, key) {
			return true
		}

		return false
	}
}

func KeyGreat[K cmp.Ordered, V any](value K) object.MapCondition[K, V] {
	return func(key K, _ V) bool {
		if key > value {
			return true
		}

		return false
	}
}

func KeyLess[K cmp.Ordered, V any](value K) object.MapCondition[K, V] {
	return func(key K, _ V) bool {
		if key < value {
			return true
		}

		return false
	}
}

func KeyRange[K cmp.Ordered, V any](i, j K) object.MapCondition[K, V] {
	return func(key K, _ V) bool {
		if key >= i && key < j {
			return true
		}

		return false
	}
}

func ValueGreat[K int, I cmp.Ordered](value I) object.MapCondition[K, I] {
	return func(index K, item I) bool {
		if item > value {
			return true
		}

		return false
	}
}

func ValueGreatOrEqual[K int, I cmp.Ordered](value I) object.MapCondition[K, I] {
	return func(index K, item I) bool {
		if item >= value {
			return true
		}

		return false
	}
}

func ValueLess[K int, I cmp.Ordered](value I) object.MapCondition[K, I] {
	return func(_ K, item I) bool {
		if item > value {
			return true
		}

		return false
	}
}

func ValueLessOrEqual[K int, I cmp.Ordered](value I) object.MapCondition[K, I] {
	return func(_ K, item I) bool {
		if item >= value {
			return true
		}

		return false
	}
}

func ValueRange[K int, I cmp.Ordered](from, to I) object.MapCondition[K, I] {
	return func(_ K, item I) bool {
		if item >= from && item < to {
			return true
		}

		return false
	}
}

func ValueEqual[K int, I comparable](items ...I) object.MapCondition[K, I] {
	return func(_ K, item I) bool {
		if slices.Contains(items, item) {
			return true
		}

		return false
	}
}

func ValueNotEqual[K int, I comparable](items ...I) object.MapCondition[K, I] {
	return func(_ K, item I) bool {
		if !slices.Contains(items, item) {
			return true
		}

		return false
	}
}

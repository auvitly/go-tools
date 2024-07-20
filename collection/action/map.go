package action

import (
	"github.com/auvitly/go-tools/collection/object"
)

func FilterMap[K comparable, V any](conditions ...object.MapCondition[K, V]) object.MapAction[K, V] {
	var del object.MapCondition[K, V] = func(key K, value V) bool {
		for i := range conditions {
			if conditions[i](key, value) {
				return true
			}
		}

		return false
	}

	return func(m object.Map[K, V]) object.Map[K, V] {
		for k, v := range m {
			if del(k, v) {
				delete(m, k)
			}
		}

		return m
	}
}

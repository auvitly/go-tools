package object

import (
	"maps"
	"slices"
)

type (
	Map[K comparable, V any]          map[K]V
	MapAction[K comparable, V any]    func(Map[K, V]) Map[K, V]
	MapCondition[K comparable, V any] func(K, V) bool
)

// MapOf - method for converting a map into object.Map.
func MapOf[K comparable, V any](m map[K]V) Map[K, V] { return m }

func (m Map[K, V]) Len() int { return len(m) }

func (m Map[K, V]) CopyTo(dst Map[K, V]) Map[K, V] { maps.Copy(dst, m); return m }

func (m Map[K, V]) CopyFrom(src Map[K, V]) Map[K, V] { maps.Copy(m, src); return m }

func (m Map[K, V]) ScanKV(keys Slice[K], values Slice[V]) Map[K, V] {
	if keys.Len() != values.Len() {
		panic("length of the keys and values does not match")
	}

	for i := range keys {
		m[keys[i]] = values[i]
	}

	return m
}

func (m Map[K, V]) ScanSlice(src Slice[V], fn func(item V) (key K)) Map[K, V] {
	for _, item := range src {
		m[fn(item)] = item
	}

	return m
}

func (m Map[K, V]) Is(conditions ...MapCondition[K, V]) (ok bool) {
	if conditions == nil {
		return false
	}

	for key, value := range m {
		for i := range conditions {
			if !conditions[i](key, value) {
				return false
			}
		}

	}

	return true
}

func (m Map[K, V]) Delete(keys ...K) Map[K, V] {
	for key := range m {
		if slices.Contains(keys, key) {
			delete(m, key)
		}
	}

	return m
}

func (m Map[K, V]) DeleteFunc(del func(K, V) bool) Map[K, V] {
	maps.DeleteFunc(m, del)

	return m
}

func (m Map[K, V]) Action(actions ...MapAction[K, V]) Map[K, V] {
	var result = m

	for i := range actions {
		result = actions[i](result)
	}

	return result
}

func (m Map[K, V]) Keys() Slice[K] {
	var keys = make(Slice[K], 0, m.Len())

	for key := range m {
		keys = keys.Append(key)
	}

	return keys
}

func (m Map[K, V]) Values() Slice[V] {
	var values = make(Slice[V], 0, m.Len())

	for key := range m {
		values = values.Append(m[key])
	}

	return values
}

func (m Map[K, V]) Clone() Map[K, V] {
	return maps.Clone(m)
}

func (m Map[K, V]) ForEach(fn func(key K, item V)) Map[K, V] {
	if fn == nil {
		return m
	}

	for key := range m {
		fn(key, m[key])
	}

	return m
}

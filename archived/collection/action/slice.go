package action

import (
	"github.com/auvitly/go-tools/collection/object"
)

func FilterSlice[E any](conditions ...object.SliceCondition[E]) object.SliceAction[E] {
	var del object.SliceCondition[E] = func(index int, item E) bool {
		for i := range conditions {
			if conditions[i](index, item) {
				return true
			}
		}

		return false
	}

	return func(s object.Slice[E]) object.Slice[E] {
		for i, v := range s {
			if del(i, v) {
				j := i
				for i++; i < len(s); i++ {
					v = s[i]
					if !del(i, v) {
						s[j] = v
						j++
					}
				}
				return s[:j]
			}
		}

		return s
	}
}

func UniqueSlice[E comparable](s object.Slice[E]) object.Slice[E] {
	var elems = make(object.Map[E, struct{}])

	for _, elem := range s {
		elems[elem] = struct{}{}
	}

	for i, v := range s {
		if _, ok := elems[v]; !ok {
			j := i
			for i++; i < len(s); i++ {
				v = s[i]
				if _, ok = elems[v]; ok {
					s[j] = v
					j++
				}
			}
			return s[:j]
		}
	}

	return s
}

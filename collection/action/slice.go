package action

import (
	"github.com/auvitly/go-tools/collection/object"
)

func DeleteElems[E any](conditions ...object.SliceCondition[E]) object.SliceAction[E] {
	var del object.SliceCondition[E] = func(index int, item E) bool {
		for i := range conditions {
			if !conditions[i](index, item) {
				return false
			}
		}

		return true
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

func FilterUnique[E comparable](s object.Slice[E]) object.Slice[E] {
	var elems = make(object.Map[E, struct{}])

	for _, elem := range s {
		elems[elem] = struct{}{}
	}

	s = make(object.Slice[E], 0, len(elems))

	return s.Append(elems.Keys()...)
}

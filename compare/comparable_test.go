package compare_test

import (
	"github.com/auvitly/go-tools/compare"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Struct struct {
	Value int
}

func (s Struct) LessThan(value compare.ComparableLess) bool {
	if b, ok := value.(Struct); ok {
		return s.Value < b.Value
	}

	return false
}

func (s Struct) Equal(value compare.ComparableEqual) bool {
	if b, ok := value.(Struct); ok {
		return s.Value == b.Value
	}

	return false
}

func (s Struct) GreaterThan(value compare.ComparableGreater) bool {
	if b, ok := value.(Struct); ok {
		return s.Value > b.Value
	}

	return false
}

func TestComparable(t *testing.T) {
	var a, b = Struct{Value: 1}, Struct{Value: 2}

	assert.True(t, compare.Less(a, b))
	assert.True(t, compare.LessOrEqual(a, b))
	assert.True(t, compare.NotEqual(a, b))
	assert.True(t, compare.NotGreater(a, b))

	assert.False(t, compare.NotLess(a, b))
	assert.False(t, compare.Equal(a, b))
	assert.False(t, compare.Greater(a, b))
	assert.False(t, compare.Equal(a, b))
	assert.False(t, compare.GreaterOrEqual(a, b))
}

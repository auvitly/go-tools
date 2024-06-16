package lab

import "testing"

// Test - unified test with establishing behavior.
type Test[I, O any] struct {
	Title       string
	Preparation []Preparation
	Input       I
	Output      O
}

// Prepare - method for setting behavior from test description.
func (test Test[I, O]) Prepare(t *testing.T, objects ...any) Test[I, O] {
	for _, object := range objects {
		for _, item := range test.Preparation {
			item.Exec(t, object)
		}
	}

	return test
}

package lab

import "testing"

// Experiment - unified test data model with preparatory actions.
type Experiment[I, O any] struct {
	Preparations Preparations
	Name         string
	Input        I
	Output       O
}

// Init - method for setting behavior from test description.
func (exp Experiment[I, O]) Init(t *testing.T, objects ...any) {
	for _, item := range exp.Preparations {
		item.Init(t, objects...)
	}
}

func (exp Experiment[I, O]) Run(t *testing.T, f func(*testing.T, Experiment[I, O])) {
	t.Helper()

	t.Run(exp.Name, func(tt *testing.T) {
		f(tt, exp)
	})
}

// Test -  unified test model. The concept is that a test consists of experiments that are run as table tests.
type Test[I, O any] []Experiment[I, O]

func (test Test[I, O]) Run(t *testing.T, f func(t *testing.T, test Experiment[I, O])) {
	t.Helper()

	for i := range test {
		test[i].Run(t, f)
	}
}

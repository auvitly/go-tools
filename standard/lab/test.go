package lab

import "testing"

// Tests - model of test scenarios.
type Tests[I, O any] []Test[I, O]

// Test - unified testcase data model with preparatory actions.
type Test[I, O any] struct {
	// Name - short name of the test script.
	Name string
	// Description - full description of the test script (optional field).
	Description string
	// In - input data object.
	In I
	// Out - output data object.
	Out O
}

// In - input data object.
type In[R any, B any] struct {
	// Request — basic model of input parameters. Can be either of a certain type,
	// there is a structure.
	// Can be filled using lab.Empty, lab.TODO or any custom type.
	Request R
	// Behavior model of behavior. Should usually be described by a function that accepts
	// as the first parameter *testing.T or *testing.B depending on the essence of the test task.
	// Can be filled using lab.Empty, lab.TODO or any custom type.
	// Example: func(t *testing.T), lab.Empty.
	Behavior B
}

// Out - output data object.
type Out[R any, E error] struct {
	// Response model of response to request execution.
	// Can be filled using lab.Empty, lab.TODO or any custom type.
	Response R
	// Error model of an error, can be either a regular error or a specific implementation.
	// Can be filled using lab.Empty or any custom error type.
	Error E
}

// TODO placeholder for base fields that should contain values.
type TODO any

// Empty is an empty placeholder for any field of the base constructs.
type Empty struct{}

// Error returns empty string.
func (Empty) Error() string { return "" }

// Run runs table tests.
func (tests Tests[I, O]) Run(t *testing.T, fn func(t *testing.T, test Test[I, O])) {
	t.Helper()

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			t.Helper()

			fn(t, test)
		})
	}
}

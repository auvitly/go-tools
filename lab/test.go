package lab

// Test - unified test data model with preparatory actions.
type Test[I, O any] struct {
	Name        string
	Description string
	In          I
	Out         O
}

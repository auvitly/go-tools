package lab

import "testing"

// Preparations - list of defines a model for preparation operations.
type Preparations []Preparation

// Preparation - defines a model for preparation operations.
type Preparation interface {
	Init(t *testing.T, args ...any)
}

func (p Preparations) Init(t *testing.T, args ...any) {
	for _, item := range p {
		item.Init(t, args...)
	}
}

// Behavior - description of dependency calls with set method.
type Behavior[D, C any] struct {
	Data   []D
	Setter func(t *testing.T, ctrl C, data D)
	fn     func(t *testing.T, ctrl C, data []D) func()
}

// Init - setting behavior on arguments.
func (d *Behavior[D, C]) Init(t *testing.T, args ...any) {
	if d.Setter == nil || args == nil {
		return
	}

	d.fn = func(t *testing.T, ctrl C, data []D) func() {
		return func() {
			for _, item := range d.Data {
				d.Setter(t, ctrl, item)
			}
		}
	}

	for _, arg := range args {
		impl, ok := arg.(C)
		if ok {
			d.fn(t, impl, d.Data)()
		}
	}
}

package lab

import "testing"

// Preparation - defines a model for preparation operations.
type Preparation interface {
	Exec(t *testing.T, object any)
}

// Behavior - description of dependency calls with set method.
type Behavior[D, C any] struct {
	Data   []D
	Setter func(t *testing.T, ctrl C, data D)
	fn     func(t *testing.T, ctrl C, data []D) func()
}

// Exec - setting behavior on controller.
func (d *Behavior[D, C]) Exec(t *testing.T, ctrl any) {
	if d.Setter == nil || ctrl == nil {
		return
	}

	d.fn = func(t *testing.T, ctrl C, data []D) func() {
		return func() {
			for _, item := range d.Data {
				d.Setter(t, ctrl, item)
			}
		}
	}

	impl, ok := ctrl.(C)
	if ok {
		d.fn(t, impl, d.Data)()
	}
}

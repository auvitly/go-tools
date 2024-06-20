package lab

import (
	"fmt"
	"testing"
)

// Preparations - list of defines a model for preparation operations.
type Preparations []Preparation

// Preparation - defines a model for preparation operations.
type Preparation interface {
	Init(t *testing.T, ctrl any)
}

func (p Preparations) Init(t *testing.T, controllers ...any) {
	for _, item := range p {
		for _, controller := range controllers {
			item.Init(t, controller)
		}
	}
}

// Behavior - description of dependency calls with set method.
type Behavior[D, C any] struct {
	Data   []D
	Setter func(t *testing.T, ctrl C, data D)
	fn     func(t *testing.T, ctrl C, data []D) func()
}

// Init - setting behavior on arguments.
func (d *Behavior[D, C]) Init(t *testing.T, ctrl any) {
	if d.Setter == nil {
		panic(fmt.Sprintf("behavior %T not contains setter", *d))
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

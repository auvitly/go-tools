package lab

import "sync"

// Call - method call model.
type Call[D, C any] struct {
	Data   D
	Setter func(ctrl C, data D)
	fn     func(ctrl C, data D) func()
	once   sync.Once
}

// Set - setting behavior on controller.
func (c *Call[D, C]) Set(ctrl any) {
	if c.Setter == nil || ctrl == nil || c.fn != nil {
		return
	}

	c.fn = func(ctrl C, data D) func() {
		return func() {
			c.Setter(ctrl, c.Data)
		}
	}

	impl, ok := ctrl.(C)
	if ok {
		c.once.Do(c.fn(impl, c.Data))
	}
}

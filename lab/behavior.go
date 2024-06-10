package lab

// Calls - method call model.
type Calls[D, C any] struct {
	Data   []D
	Setter func(ctrl C, data D)
	fn     func(ctrl C, data []D) func()
}

// Set - setting behavior on controller.
func (c *Calls[D, C]) Set(ctrl any) {
	if c.Setter == nil || ctrl == nil {
		return
	}

	c.fn = func(ctrl C, data []D) func() {
		return func() {
			for _, item := range c.Data {
				c.Setter(ctrl, item)
			}
		}
	}

	impl, ok := ctrl.(C)
	if ok {
		c.fn(impl, c.Data)()
	}
}

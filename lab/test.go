package lab

// Test - unified test with establishing behavior.
type Test[R Request, E Expect] struct {
	// Title - allows you to set a short title that can be easily found when needed.
	Title string
	// Request - request parameters for the test.
	Request R
	// Behavior - behavior system data.
	Behavior []Behavior
	// Expected - result of program execution.
	Expected E
}

// ApplyBehavior - method for setting behavior from test description.
func (t Test[R, E]) ApplyBehavior(ctrl ...any) Test[R, E] {
	for _, _ctrl := range ctrl {
		for _, item := range t.Behavior {
			item.Set(_ctrl)
		}
	}

	return t
}

// Interface that is used to implement a set of models within a package
type (
	Behavior interface{ Set(ctrl any) }
	Request  interface{ implRequestData() }
	Expect   interface{ implExpectData() }
)

// Payload - unified payload model.
type Payload[P any] struct {
	Payload P
}

// Result - unified results model format.
type Result[P any, E error] struct {
	Payload P
	Error   E
}

// Error - payload-free unified result model format.
type Error[E error] struct {
	Error E
}

func (Payload[P]) implRequestData()  {}
func (Payload[P]) implExpectData()   {}
func (Result[P, E]) implExpectData() {}
func (Error[E]) implExpectData()     {}

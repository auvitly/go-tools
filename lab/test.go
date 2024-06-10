package lab

// SimpleTest - a simplified form of the unified test without establishing behavior.
type SimpleTest[R Request, E Expect] struct {
	// Title - allows you to set a short title that can be easily found when needed.
	Title string
	// Request - request parameters for the test.
	Request R
	// Expected - result of program execution.
	Expected E
}

// Test - unified test format.
type Test[R Request, B Behavior, E Expect] struct {
	// Title - allows you to set a short title that can be easily found when needed.
	Title string
	// Request - request parameters for the test.
	Request R
	// Behaviour - behavior system data.
	Behavior B
	// Expected - result of program execution.
	Expected E
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

func (Test[P, B, E]) implRequestData() {}
func (Payload[P]) implRequestData()    {}
func (Payload[P]) implExpectData()     {}
func (Result[P, E]) implExpectData()   {}
func (Error[E]) implExpectData()       {}

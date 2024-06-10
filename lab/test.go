package lab

// Test - unified test format.
type Test[R RequestData, E ExpectData] struct {
	// Title - allows you to set a short title that can be easily found when needed.
	Title string
	// Request - request parameters for the test.
	Request R
	// Expected - result of program execution.
	Expected E
}

// Tester - an interface that is used to implement a set of models within a package.
type Tester interface {
	implTester()
}

// Payload - unified payload model.
type Payload[P any] struct {
	Payload P
}

// RequestData - an interface that is used to implement a set of models within a package.
type RequestData interface {
	implRequestData()
}

// Behavior - an adapter model for your implementation of parameters.
type Behavior[P, B any] struct {
	Payload  P
	Behavior B
}

// ExpectData - an interface that is used to implement a set of models within a package.
type ExpectData interface {
	implExpectData()
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

func (Test[R, E]) implTester()          {}
func (Behavior[D, B]) implRequestData() {}
func (Payload[D]) implRequestData()     {}
func (Payload[D]) implExpectData()      {}
func (Result[D, E]) implExpectData()    {}
func (Error[E]) implExpectData()        {}

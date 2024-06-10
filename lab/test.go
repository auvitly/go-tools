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

func (Test[R, E]) implTester() {}

// Data - an adapter model for your implementation of parameters.
type Data[D any] struct {
	Data D
}

func (Data[D]) implExpectData()  {}
func (Data[D]) implRequestData() {}

// RequestData - an interface that is used to implement a set of models within a package.
type RequestData interface {
	implRequestData()
}

// DataWithBehavior - an adapter model for your implementation of parameters.
type DataWithBehavior[D, B any] struct {
	Data[D]
	Behavior B
}

func (DataWithBehavior[D, B]) implRequestData() {}

// ExpectData - an interface that is used to implement a set of models within a package.
type ExpectData interface {
	implExpectData()
}

type Result[D any, E error] struct {
	Data[D]
	Error E
}

type Error[E error] struct {
	Error E
}

func (Result[D, E]) implExpectData() {}
func (Error[E]) implExpectData()     {}

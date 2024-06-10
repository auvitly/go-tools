package lab

// Empty - .
type Empty struct{}

func (Empty) Error() string    { return "" }
func (Empty) implRequestData() {}
func (Empty) implExpectData()  {}
func (Empty) Set(ctrl any)     {}

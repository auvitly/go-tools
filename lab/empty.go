package lab

// Empty - .
type Empty struct{}

func (Empty) Error() string    { return "" }
func (Empty) Set(ctrl any)     {}
func (Empty) implRequestData() {}
func (Empty) implExpectData()  {}

package lab

type Empty struct{}

func (Empty) Error() string    { return "" }
func (Empty) implRequestData() {}
func (Empty) implExpectData()  {}

package async

var _ch = makeClosedChan[struct{}]

func makeClosedChan[T any]() chan T {
	var ch = make(chan T)

	close(ch)

	return ch
}

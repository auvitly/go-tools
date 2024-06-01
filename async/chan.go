package async

var _ch = func() chan struct{} {
	var ch = make(chan struct{})

	close(ch)

	return ch
}

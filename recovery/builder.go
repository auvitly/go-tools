package recovery

import "github.com/auvitly/go-tools/stderrs"

type Handler func(msg any)

type Builder struct {
	handlers []Handler
	target   *error
	stderr   **stderrs.Error
	message  string
}

var (
	_handlers []Handler
	_message  = "internal server error: unhandled exception"
)

func (b Builder) SetMessage(message string) Builder {
	var dst = b.copy()

	dst.message = message

	return dst
}

func (b Builder) OnError(err *error) Builder {
	var dst = b.copy()

	dst.target = err
	dst.stderr = nil

	return dst
}

func (b Builder) On(err **stderrs.Error) Builder {
	var dst = b.copy()

	dst.target = nil
	dst.stderr = err

	return dst
}

func (b Builder) WithHandlers(handlers ...Handler) Builder {
	var dst = b.copy()

	dst.handlers = append(dst.handlers, handlers...)

	return dst
}

func (b Builder) copy() Builder {
	var list = make([]Handler, 0, len(b.handlers))

	list = append(list, b.handlers...)

	return Builder{
		handlers: list,
		target:   b.target,
		stderr:   b.stderr,
		message:  b.message,
	}
}

func (b Builder) Do() {
	if msg := recover(); msg != nil {
		b.recovery(msg)
	}
}

func (b Builder) recovery(msg any) {
	var (
		errs []error
		use  = func(handler Handler) {
			defer func() {
				if sub := recover(); sub != nil {
					errs = append(errs, stderrs.Panic.
						SetMessage(b.message).
						WithField("panic", sub),
					)
				}
			}()

			handler(msg)
		}
	)

	if len(b.message) == 0 {
		b.message = _message
	}

	for _, handler := range _handlers {
		use(handler)
	}

	for _, handler := range b.handlers {
		use(handler)
	}

	switch {
	case b.target != nil:
		var std = stderrs.Panic.
			SetMessage(b.message).
			EmbedErrors(errs...).
			WithField("panic", msg)

		if *b.target != nil {
			std = std.EmbedErrors(*b.target)
		}

		*b.target = std
	case b.stderr != nil:
		var std = stderrs.Panic.
			SetMessage(b.message).
			EmbedErrors(errs...).
			WithField("panic", msg)

		if *b.stderr != nil {
			std = std.EmbedErrors(*b.stderr)
		}

		*b.stderr = std
	}
}

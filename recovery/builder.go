package recovery

import (
	"fmt"
	"github.com/auvitly/go-tools/stderrs"
	"runtime/debug"
	"slices"
)

// Handler - user panic handler.
type Handler func(msg any) (err error)

// AsyncHandler - async user panic handler.
type AsyncHandler func(msg any)

// Builder - panic builder.
type Builder struct {
	syncHandlers  []Handler
	asyncHandlers []AsyncHandler
	target        *error
	stderr        **stderrs.Error
	message       string
	enriched      bool
}

// SetMessage - set message for standard error.
func (b Builder) SetMessage(message string) Builder {
	var dst = b.copy()

	dst.message = message

	return dst
}

// OnError - perform error enrichment.
func (b Builder) OnError(err *error) Builder {
	var dst = b.copy()

	dst.target = err
	dst.stderr = nil

	return dst
}

// On - perform standard error enrichment.
func (b Builder) On(err **stderrs.Error) Builder {
	var dst = b.copy()

	dst.target = nil
	dst.stderr = err

	return dst
}

// WithHandlers - add exception handler.
func (b Builder) WithHandlers(handlers ...Handler) Builder {
	var dst = b.copy()

	dst.syncHandlers = append(dst.syncHandlers, handlers...)

	return dst
}

// WithAsyncHandlers - add async exception handler.
func (b Builder) WithAsyncHandlers(handlers ...AsyncHandler) Builder {
	var dst = b.copy()

	dst.asyncHandlers = append(dst.asyncHandlers, handlers...)

	return dst
}

// WithoutHandlers - allows you to reset all handlers for the selected call.
func (b Builder) WithoutHandlers() Builder {
	var dst = b.copy()

	dst.asyncHandlers = nil
	dst.syncHandlers = nil

	return dst
}

func (b Builder) copy() Builder {
	return Builder{
		syncHandlers:  slices.Clone(b.syncHandlers),
		asyncHandlers: slices.Clone(b.asyncHandlers),
		target:        b.target,
		stderr:        b.stderr,
		message:       b.message,
		enriched:      true,
	}
}

// Do - perform panic processing with context. Called exclusively via defer.
func (b Builder) Do() {
	if msg := recover(); msg != nil {
		b.recovery(msg)
	}
}

func (b Builder) useSync(
	msg any,
	errs *[]error,
	handler Handler,
) {
	var err error

	defer func() {
		if sub := recover(); sub != nil {
			var std = stderrs.Panic.
				WithField("panic", fmt.Sprintf("%s", sub)).
				WithField("stack", string(debug.Stack()))

			if err != nil {
				std = std.EmbedErrors(err)
			}

			*errs = append(*errs, std)

			return
		}

		if err != nil {
			*errs = append(*errs, err)
		}
	}()

	err = handler(msg)
}

func (b Builder) recovery(msg any) {
	switch {
	case !b.enriched:
		return
	case (b.stderr != nil || b.target != nil) && len(b.asyncHandlers) == 0 && len(b.syncHandlers) == 0:
		b.setError(nil, msg)

		return
	}

	var errs []error

	if len(b.message) == 0 {
		b.message = _message
	}

	if len(b.asyncHandlers) != 0 {
		for _, handler := range b.asyncHandlers {
			go handler(msg)
		}
	}

	if len(b.syncHandlers) != 0 {
		for _, handler := range b.syncHandlers {
			b.useSync(msg, &errs, handler)
		}
	}

	b.setError(errs, msg)

	return
}

func (b Builder) setError(errs []error, msg any) {
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

func (b Builder) handle(
	msg any,
	errs *[]error,
) {

}

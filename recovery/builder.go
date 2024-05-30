package recovery

import (
	"context"
	"fmt"
	"github.com/auvitly/go-tools/stderrs"
	"runtime/debug"
	"slices"
	"sync"
)

// Handler - sync user panic handler.
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
		b.recovery(context.Background(), msg)
	}
}

// DoContext - perform panic processing with context. Called exclusively via defer.
func (b Builder) DoContext(ctx context.Context) {
	if msg := recover(); msg != nil {
		b.recovery(ctx, msg)
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

func (b Builder) useAsync(
	ctx context.Context,
	msg any,
	mu *sync.Mutex,
	errs *[]error,
	handler AsyncHandler,
) {
	defer func() {
		var sub = recover()

		if sub == nil {
			return
		}

		var std = stderrs.Panic.
			WithField("panic", fmt.Sprintf("%s", sub)).
			WithField("stack", string(debug.Stack()))
		select {
		case <-ctx.Done():
			return
		default:
			mu.Lock()
			*errs = append(*errs, std)
			mu.Unlock()
		}
	}()

	handler(msg)
}

func (b Builder) recovery(ctx context.Context, msg any) {
	switch {
	case !b.enriched:
		return
	case (b.stderr != nil || b.target != nil) && len(b.asyncHandlers) == 0 && len(b.syncHandlers) == 0:
		b.setError(nil, msg)

		return
	}

	var (
		errs []error
		ch   chan struct{}
	)

	if len(b.message) == 0 {
		b.message = _message
	}

	if len(b.asyncHandlers) != 0 || len(b.syncHandlers) != 0 {
		b.handle(ctx, msg, &ch, &errs)
	}

	if len(b.asyncHandlers) == 0 {
		b.setError(errs, msg)

		return
	}

	for {
		select {
		case <-ctx.Done():
			b.setError(errs, msg)

			return
		case <-ch:
			b.setError(errs, msg)

			return
		}
	}
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
	ctx context.Context,
	msg any,
	ch *chan struct{},
	errs *[]error,
) {
	var mu sync.Mutex

	if len(b.asyncHandlers) != 0 {
		var wg sync.WaitGroup

		*ch = make(chan struct{})

		for _, handler := range b.asyncHandlers {
			wg.Add(1)

			go func(handler AsyncHandler) {
				defer wg.Done()

				b.useAsync(ctx, msg, &mu, errs, handler)
			}(handler)
		}

		go func() {
			wg.Wait()
			*ch <- struct{}{}
		}()
	}

	mu.Lock()
	defer mu.Unlock()

	for _, handler := range b.syncHandlers {
		b.useSync(msg, errs, handler)
	}
}

package recovery

import (
	"context"
	"fmt"
	"github.com/auvitly/go-tools/stderrs"
	"runtime/debug"
	"sync"
)

// SyncHandler - sync user panic handler.
type SyncHandler func(msg any) (err error)

// AsyncHandler - async user panic handler.
type AsyncHandler func(msg any)

// Builder - panic builder.
type Builder struct {
	syncHandlers  []SyncHandler
	asyncHandlers []AsyncHandler
	target        *error
	stderr        **stderrs.Error
	message       string
	global        bool
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

// WithSyncHandlers - add exception handler.
func (b Builder) WithSyncHandlers(handlers ...SyncHandler) Builder {
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
	var (
		syncHandlers  = make([]SyncHandler, 0, len(b.syncHandlers))
		asyncHandlers = make([]AsyncHandler, 0, len(b.asyncHandlers))
	)

	return Builder{
		syncHandlers:  append(syncHandlers, b.syncHandlers...),
		asyncHandlers: append(asyncHandlers, b.asyncHandlers...),
		target:        b.target,
		stderr:        b.stderr,
		message:       b.message,
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
	mu *sync.Mutex,
	errs *[]error,
	handler SyncHandler,
) {
	var err error

	defer func() {
		var sub = recover()

		if sub == nil {
			return
		}

		var std = stderrs.Panic.
			WithField("panic", fmt.Sprintf("%s", sub)).
			WithField("stack", string(debug.Stack()))

		if err != nil {
			std = std.EmbedErrors(err)
		}

		mu.Lock()
		*errs = append(*errs, std)
		mu.Unlock()
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
	var (
		errs []error
		wg   sync.WaitGroup
		mu   sync.Mutex
		ch   = make(chan struct{})
	)

	if len(b.message) == 0 {
		b.message = _message
	}

	b.handle(ctx, msg, &errs, &wg, &mu)

	go func() {
		wg.Wait()
		ch <- struct{}{}
	}()

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
	errs *[]error,
	wg *sync.WaitGroup,
	mu *sync.Mutex,
) {
	for _, handler := range b.asyncHandlers {
		wg.Add(1)

		go func(handler AsyncHandler) {
			defer wg.Done()

			b.useAsync(ctx, msg, mu, errs, handler)
		}(handler)
	}

	for _, handler := range b.syncHandlers {
		b.useSync(msg, mu, errs, handler)
	}
}

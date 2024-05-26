package recovery

import (
	"context"
	"github.com/auvitly/go-tools/stderrs"
	"sync"
)

type Handler func(ctx context.Context, msg any) error

type Builder struct {
	syncHandlers  []Handler
	asyncHandlers []Handler
	target        *error
	stderr        **stderrs.Error
	message       string
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
func (b Builder) WithAsyncHandlers(handlers ...Handler) Builder {
	var dst = b.copy()

	dst.syncHandlers = append(dst.syncHandlers, handlers...)

	return dst
}

func (b Builder) copy() Builder {
	var (
		syncHandlers  = make([]Handler, 0, len(b.syncHandlers))
		asyncHandlers = make([]Handler, 0, len(b.asyncHandlers))
	)

	return Builder{
		syncHandlers:  append(syncHandlers, b.syncHandlers...),
		asyncHandlers: append(asyncHandlers, b.asyncHandlers...),
		target:        b.target,
		stderr:        b.stderr,
		message:       b.message,
	}
}

// Do - perform panic processing. Called exclusively via defer.
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

func (b Builder) use(
	ctx context.Context,
	msg any,
	mu *sync.Mutex,
	errs *[]error,
	handler Handler,
) {
	var err error

	defer func() {
		if sub := recover(); sub != nil {
			var std = stderrs.Panic.
				SetMessage(b.message).
				WithField("panic", sub)

			if err != nil {
				std = std.EmbedErrors(err)
			}

			mu.Lock()
			*errs = append(*errs, std)
			mu.Unlock()

			return
		}

		if err != nil {
			mu.Lock()
			*errs = append(*errs, err)
			mu.Unlock()
		}
	}()

	err = handler(ctx, msg)
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
	for _, handler := range _asyncHandlers {
		wg.Add(1)

		go func(handler Handler) {
			defer wg.Done()

			b.use(ctx, msg, mu, errs, handler)
		}(handler)
	}

	for _, handler := range b.asyncHandlers {
		wg.Add(1)

		go func(handler Handler) {
			defer wg.Done()

			b.use(ctx, msg, mu, errs, handler)
		}(handler)
	}

	wg.Add(1)

	go func() {
		defer wg.Done()

		for _, handler := range _syncHandlers {
			b.use(ctx, msg, mu, errs, handler)
		}

		for _, handler := range b.syncHandlers {
			b.use(ctx, msg, mu, errs, handler)
		}
	}()
}

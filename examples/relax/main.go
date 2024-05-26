package main

import (
	"context"
	"fmt"
	"github.com/auvitly/go-tools/recovery"
	"github.com/auvitly/go-tools/stderrs"
	"io/fs"
	"log/slog"
	"time"
)

// log - logs the panic message.
func log(ctx context.Context, msg any) error {
	slog.WarnContext(ctx, fmt.Sprintf("%s", msg))

	return nil
}

// errClosed - the handler failed with an error fs.ErrClosed.
func errClosed(ctx context.Context, msg any) error {
	return fs.ErrClosed
}

// wrap - the method demonstrates the ability to wrap a function to pass arguments to a handler.
func wrap(text string) func(context.Context, any) error {
	return func(_ context.Context, msg any) error {
		return stderrs.Internal.SetMessage("%s: %s", text, msg)
	}
}

// itsPanic - panic exit handler.
func itsPanic(context.Context, any) error {
	var i *int

	_ = *i

	return nil
}

// tooLate - the time required to execute this method exceeds the context timeout.
func tooLate(ctx context.Context, _ any) error {
	time.Sleep(2 * time.Second)

	slog.InfoContext(ctx, "I won't get to the output because the context ended early")

	return nil
}

// onTime - a method that manages to asynchronously output a message to the log.
func onTime(ctx context.Context, _ any) error {
	slog.InfoContext(ctx, "I slipped through")

	return nil
}

// method - panic occurs and calls the recovery package handler.
func method(ctx context.Context) (err *stderrs.Error) {
	defer recovery.On(&err).Do(ctx)

	panic("I'm dropping the app now! Be afraid!")

	return nil
}

func main() {
	recovery.RegistryHandlers(
		log,
		errClosed,
		wrap("message"),
		itsPanic,
	)

	recovery.RegistryAsyncHandlers(
		tooLate,
		onTime,
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := method(ctx); err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("Panic did not overtake us! We received an error: %s", err))
	}
}

/* OUT:
> [ INFO ]  I managed to get to the conclusion
> [ WARN ]  I'm dropping the app now! Be afraid!
> [ ERROR ] Panic did not overtake us! We received an error: {
	"code": "panic",
	"message": "internal server error: unhandled exception",
	"fields": {
		"panic":"I'm dropping the app now! Be afraid!"
	},
	"embed": [
		file already closed
		{
			"code": "internal",
			"message": "message: I'm dropping the app now! Be afraid!"
		}
		{
			"code": "panic",
			"message": "internal server error: unhandled exception",
			"fields": {
				"panic":"invalid memory address or nil pointer dereference",
				"stack":"goroutine 8 [running]:
					runtime/debug.Stack()
						C:/Program Files/Go/src/runtime/debug/stack.go:24 +0x5e
							github.com/auvitly/go-tools/recovery.Builder.use.func1()
						F:/Work/projects/git/auvitly/go-tools/recovery/builder.go:103 +0x433
							panic({0x9141a0?, 0xb9f920?})
						C:/Program Files/Go/src/runtime/panic.go:914 +0x21f
							main.itsPanic({0x9cfcb0, 0xc00007a0e0}, {0x8fcc80, 0x9cd8b0})
						F:/Work/projects/git/auvitly/go-tools/examples/relax/main.go:36 +0x2
							github.com/auvitly/go-tools/recovery.Builder.use({{0xc08360, 0x0, 0x0},{0xc08360, 0x0, 0x0}, 0x0, 0xc000044058, {0x965fa7, 0x2a}}, ...)
						F:/Work/projects/git/auvitly/go-tools/recovery/builder.go:123 +0xee
`							github.com/auvitly/go-tools/recovery.Builder.handle.func3()
						F:/Work/projects/git/auvitly/go-tools/recovery/builder.go:219 +0x11f
							created by github.com/auvitly/go-tools/recovery.Builder.handle in goroutine 1
						F:/Work/projects/git/auvitly/go-tools/recovery/builder.go:215 +0x63d\n"
			}
		}
	]
}
*/

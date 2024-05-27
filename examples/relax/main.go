package main

import (
	"context"
	"fmt"
	"github.com/auvitly/go-tools/recovery"
	"github.com/auvitly/go-tools/stderrs"
	"log/slog"
	"time"
)

// log - logs the panic message.
func log(ctx context.Context, msg any) {
	slog.InfoContext(ctx, "log", "message", msg)
}

// wrap - the method demonstrates the ability to wrap a function to pass arguments to a handler.
func wrap(text string) func(context.Context, any) {
	return func(ctx context.Context, msg any) {
		slog.InfoContext(ctx, "wrap",
			"text", text,
			"message", msg,
		)
	}
}

// itsPanic - panic exit handler.
func itsPanic(context.Context, any) {
	var i *int

	_ = *i
}

// tooLate - the time required to execute this method exceeds the context timeout.
func tooLate(ctx context.Context, _ any) {
	time.Sleep(2 * time.Second)

	slog.InfoContext(ctx, "tooLate")
}

// tooLatePanic - panic exit handler, the time required to execute this method exceeds the context timeout.
func tooLatePanic(ctx context.Context, _ any) {
	time.Sleep(2 * time.Second)

	panic("tooLatePanic")
}

// onTime - a method that manages to asynchronously output a message to the log.
func onTime(ctx context.Context, _ any) {
	slog.InfoContext(ctx, "I slipped through")
}

// method - panic occurs and calls the recovery package handler.
func method(ctx context.Context) (err *stderrs.Error) {
	defer recovery.On(&err).Do(ctx)

	err = stderrs.Unknown.SetMessage("Successfully assigned the error! Wow!")

	panic("I'm dropping the app now! Be afraid!")

	return nil
}

func main() {
	recovery.RegistryHandlers(
		log,
		wrap("message"),
		itsPanic,
	)

	recovery.RegistryAsyncHandlers(
		tooLate,
		onTime,
		tooLatePanic,
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := method(ctx); err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("Panic did not overtake us! We received an error: %s", err))
	}

	time.Sleep(2 * time.Second)
}

/* OUT:
2024/05/27 23:26:04 INFO I slipped through
2024/05/27 23:26:04 INFO log message="I'm dropping the app now! Be afraid!"
2024/05/27 23:26:04 INFO wrap text="message" message="I'm dropping the app now! Be afraid!"
2024/05/27 23:26:05 ERROR Panic did not overtake us! We received an error: {
	"code": "panic",
	"message": "internal server error: unhandled exception",
	"fields": {
		"panic":"I'm dropping the app now! Be afraid!",
		"uuid":"179e3034-270e-48d3-9459-d83cf89545a8",
	},
	"embed": [
		{
			"code": "panic",
			"fields": {
				"panic":"runtime error: invalid memory address or nil pointer dereference",
				"uuid":"179e3034-270e-48d3-9459-d83cf89545a8",
				"stack":"goroutine 9 [running]:
					runtime/debug.Stack()
						C:/Program Files/Go/src/runtime/debug/stack.go:24 +0x5e
					github.com/auvitly/go-tools/recovery.Builder.use.func1()
						F:/Work/projects/git/auvitly/go-tools/recovery/builder.go:111 +0x2be
					panic({0xd953c0?, 0x10229c0?})
						C:/Program Files/Go/src/runtime/panic.go:914 +0x21f
					main.itsPanic({0xe516f0, 0xc00007a0e0}, {0xd7dde0, 0xe4f2c0})
						F:/Work/projects/git/auvitly/go-tools/examples/relax/main.go:31 +0x2
					github.com/auvitly/go-tools/recovery.Builder.use({{0xb, 0x33, 0x8e, 0x7d, 0xda, 0xa7, 0x48, 0xf5, 0xb3, 0x96, ...}, ...}, ...)
						F:/Work/projects/git/auvitly/go-tools/recovery/builder.go:124 +0xf1
					github.com/auvitly/go-tools/recovery.Builder.handle.func3()
						F:/Work/projects/git/auvitly/go-tools/recovery/builder.go:224 +0x125
					created by github.com/auvitly/go-tools/recovery.Builder.handle in goroutine 1
						F:/Work/projects/git/auvitly/go-tools/recovery/builder.go:220 +0x64b",
			}
		},
		{
			"code": "unknown",
			"message": "Successfully assigned the error! Wow!"
		}
	]
}
2024/05/27 23:26:06 ERROR panic detected when executing handler after interceptor context ends:
{
	"code": "panic",
	"fields": {
		"panic":"tooLatePanic",
		"uuid":"179e3034-270e-48d3-9459-d83cf89545a8",
		"stack":"goroutine 8 [running]:
			runtime/debug.Stack()
				C:/ProgramFiles/Go/src/runtime/debug/stack.go:24 +0x5e
			github.com/auvitly/go-tools/recovery.Builder.use.func1()
				F:/Work/projects/git/auvitly/go-tools/recovery/builder.go:111 +0x2be
			panic({0x5adde0?, 0x67f2a0?})
				C:/Program Files/Go/src/runtime/panic.go:914 +0x21f
			main.tooLatePanic({0x0?, 0x0?}, {0x0?, 0x0?})
				F:/Work/projects/git/auvitly/go-tools/examples/relax/main.go:45 +0x2b
			github.com/auvitly/go-tools/recovery.Builder.use({{0x17, 0x9e, 0x30, 0x34, 0x27, 0xe, 0x48, 0xd3, 0x94, 0x59, ...}, ...}, ...)
				F:/Work/projects/git/auvitly/go-tools/recovery/builder.go:128 +0xf1
			github.com/auvitly/go-tools/recovery.Builder.handle.func1(0x0?)
				F:/Work/projects/git/auvitly/go-tools/recovery/builder.go:208 +0xc5
			created by github.com/auvitly/go-tools/recovery.Builder.handle in goroutine 1
				F:/Work/projects/git/auvitly/go-tools/recovery/builder.go:205 +0x89"
	}
}
*/

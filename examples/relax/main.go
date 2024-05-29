package main

import (
	"context"
	"fmt"
	"github.com/auvitly/go-tools/recovery"
	"github.com/auvitly/go-tools/stderrs"
	"log/slog"
	"time"
)

// syncHandler - logs the panic message.
func syncHandler(ctx context.Context) func(msg any) error {
	return func(msg any) error {
		slog.InfoContext(ctx, "log", "message", msg)

		return nil
	}
}

// syncPanicHandler - panic exit handler.
func syncPanicHandler(any) error {
	panic("syncPanicHandler")

	return nil
}

// asyncOnTime - a method that manages to asynchronously output a message to the log.
func asyncOnTime(ctx context.Context) func(any) {
	return func(any) {
		slog.InfoContext(ctx, "I slipped through")
	}
}

// asyncOnTimePanic - a method that manages to asynchronously output a message to the log.
func asyncOnTimePanic(ctx context.Context) func(any) {
	return func(any) {
		panic("asyncOnTimePanic")
	}
}

// asyncTooLateHandler - the time required to execute this method exceeds the context timeout.
func asyncTooLateHandler(ctx context.Context) func(_ any) {
	return func(any) {
		time.Sleep(2 * time.Second)

		slog.InfoContext(ctx, "tooLate")
	}
}

// asyncTooLatePanicHandler - panic exit handler, the time required to execute this method exceeds the context timeout.
func asyncTooLatePanicHandler(_ any) {
	defer recovery.WithoutHandlers().WithSyncHandlers(func(msg any) (err error) {
		slog.Error("asyncTooLatePanicHandler: recovery", "msg", msg)

		return nil
	}).Do()

	time.Sleep(2 * time.Second)

	panic("asyncTooLatePanicHandler")
}

// method - panic occurs and calls the recovery package handler.
func method(ctx context.Context) (err *stderrs.Error) {
	defer recovery.On(&err).DoContext(ctx)

	err = stderrs.Unknown.SetMessage("Successfully assigned the error! Wow!")

	panic("I'm dropping the app now! Be afraid!")

	return nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	recovery.RegistrySyncHandlers(
		syncHandler(ctx),
		syncPanicHandler,
	)

	recovery.RegistryAsyncHandlers(
		asyncOnTime(ctx),
		asyncOnTimePanic(ctx),
		asyncTooLateHandler(ctx),
		asyncTooLatePanicHandler,
	)

	if err := method(ctx); err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("Panic did not overtake us! We received an error: %s", err))
	}

	time.Sleep(10 * time.Second)
}

/* OUT:
2024/05/28 22:49:03 INFO log message="I'm dropping the app now! Be afraid!"
2024/05/28 22:49:03 INFO I slipped through
2024/05/28 22:49:04 ERROR Panic did not overtake us! We received an error: {
	"code": "panic",
	"message": "internal server error: unhandled exception",
	"fields": {
		"panic":"I'm dropping the app now! Be afraid!"
	},
	"embed": [
		{
			"code": "panic",
			"fields": {
				"panic":"asyncOnTimePanic",
				"stack":"goroutine 7 [running]:
					runtime/debug.Stack()\n\tC:/Program Files/Go/src/runtime/debug/stack.go:24 +0x5e
					github.com/auvitly/go-tools/recovery.Builder.useAsync.func1()
						F:/Work/projects/git/auvitly/go-tools/recovery/builder.go:149 +0x28a
					panic({0x64cd80?, 0x71db70?})
						C:/Program Files/Go/src/runtime/panic.go:914 +0x21f
					main.main.asyncOnTimePanic.func3({0x0, 0x0})
						F:/Work/projects/git/auvitly/go-tools/examples/relax/main.go:38 +0x25
					github.com/auvitly/go-tools/recovery.Builder.useAsync({{0x959360, 0x0, 0x0}, {0x959360, 0x0, 0x0}, 0x0, 0xc000044060, {0x6b61ac, 0x2a}}, ...)
						F:/Work/projects/git/auvitly/go-tools/recovery/builder.go:160 +0x8d
					github.com/auvitly/go-tools/recovery.Builder.handle.func1(0x0?)
						F:/Work/projects/git/auvitly/go-tools/recovery/builder.go:236 +0xc5
					created by github.com/auvitly/go-tools/recovery.Builder.handle in goroutine 1
						F:/Work/projects/git/auvitly/go-tools/recovery/builder.go:233 +0x92"
			}
		},
		{
			"code": "panic",
			"fields": {
				"panic":"syncPanicHandler",
				"stack":"goroutine 1 [running]:
					runtime/debug.Stack()
						C:/Program Files/Go/src/runtime/debug/stack.go:24 +0x5e
					github.com/auvitly/go-tools/recovery.Builder.useSync.func1()
						F:/Work/projects/git/auvitly/go-tools/recovery/builder.go:119 +0x28a
					panic({0x64cd80?, 0x71db40?})
						C:/Program Files/Go/src/runtime/panic.go:914 +0x21f
					main.syncPanicHandler({0x64cd80, 0x71db60})
						F:/Work/projects/git/auvitly/go-tools/examples/relax/main.go:23 +0x25
					github.com/auvitly/go-tools/recovery.Builder.useSync({{0x959360, 0x0, 0x0}, {0x959360, 0x0, 0x0}, 0x0, 0xc000044060, {0x6b61ac, 0x2a}}, ...)
						F:/Work/projects/git/auvitly/go-tools/recovery/builder.go:130 +0x77
					github.com/auvitly/go-tools/recovery.Builder.handle({{0x959360, 0x0, 0x0}, {0x959360, 0x0, 0x0}, 0x0, 0xc000044060, {0x6b61ac, 0x2a}}, ...)
						F:/Work/projects/git/auvitly/go-tools/recovery/builder.go:251 +0x58e
					github.com/auvitly/go-tools/recovery.Builder.recovery({{0x959360, 0x0, 0x0}, {0x959360, 0x0, 0x0}, 0x0, 0xc000044060, {0x6b61ac, 0x2a}}, ...)
						F:/Work/projects/git/auvitly/go-tools/recovery/builder.go:175 +0x11b
					github.com/auvitly/go-tools/recovery.Builder.DoContext({{0x959360, 0x0, 0x0}, {0x959360, 0x0, 0x0}, 0x0, 0xc000044060, {0x0, 0x0}}, ...)
						F:/Work/projects/git/auvitly/go-tools/recovery/builder.go:98 +0x6e
					panic({0x64cd80?, 0x71db60?})
						C:/Program Files/Go/src/runtime/panic.go:920 +0x270
					main.method({0x71ff70, 0xc00007a0e0})
						F:/Work/projects/git/auvitly/go-tools/examples/relax/main.go:64 +0x405
					main.main()
						F:/Work/projects/git/auvitly/go-tools/examples/relax/main.go:85 +0x328",
			}
		},
		{
			"code": "unknown",
			"message": "Successfully assigned the error! Wow!"
		}
	]
}
2024/05/28 23:04:33 INFO tooLate
2024/05/29 19:56:08 ERROR asyncTooLatePanicHandler: recovery msg=asyncTooLatePanicHandler
*/

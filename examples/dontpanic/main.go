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

func logHandler(ctx context.Context, msg any) error {
	slog.WarnContext(ctx, fmt.Sprintf("%s", msg))

	return nil
}

func errorHandler(ctx context.Context, msg any) error {
	return fs.ErrClosed
}

func wrapHandler(data string) func(ctx context.Context, msg any) error {
	return func(_ context.Context, _ any) error {
		return stderrs.Internal.SetMessage("My data error: %s", data)
	}
}

func exceedingTimeoutHandler(ctx context.Context, _ any) error {
	time.Sleep(2 * time.Second)

	slog.InfoContext(ctx, "I won't get to the output because the context ended early")

	return nil
}

func asyncHandler(ctx context.Context, _ any) error {
	slog.InfoContext(ctx, "I slipped through")

	return nil
}

func onStart(ctx context.Context) (err *stderrs.Error) {
	defer recovery.On(&err).DoContext(ctx)

	panic("I'm dropping the app now! Be afraid!")

	return nil
}

func main() {
	recovery.RegistryHandlers(
		logHandler,
		errorHandler,
		wrapHandler("message"),
	)

	recovery.RegistryAsyncHandlers(
		exceedingTimeoutHandler,
		asyncHandler,
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := onStart(ctx); err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("Panic did not overtake us! We received an error: %s", err))
	}
}

/* OUT:
> [ INFO ]  I managed to get to the conclusion
> [ WARN ]  I'm dropping the app now! Be afraid!
> [ ERROR ] Panic did not overtake us! We received an error: {"code": "panic", "message": "internal server error: unhandled exception", "fields": {"panic":"I'm dropping the app now! Be afraid!"}, "embed": [file already closed
{"code": "internal", "message": "My data error: message"}]}
*/

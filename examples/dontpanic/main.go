package main

import (
	"context"
	"fmt"
	"github.com/auvitly/go-tools/recovery"
	"github.com/auvitly/go-tools/stderrs"
	"log/slog"
	"time"
)

func logHandler(ctx context.Context, msg any) error {
	slog.WarnContext(ctx, fmt.Sprintf("%s", msg))

	return nil
}

func errorHandler(ctx context.Context, msg any) error {
	return stderrs.Internal.SetMessage("It's error!")
}

func onStart(ctx context.Context) (err *stderrs.Error) {
	defer recovery.On(&err).DoContext(ctx)

	panic("I'm dropping the app now! Be afraid!")

	return nil
}

func main() {
	recovery.RegistryHandlers(logHandler, errorHandler)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := onStart(ctx); err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("Panic did not overtake us! We received an error: %s", err))
	}
}

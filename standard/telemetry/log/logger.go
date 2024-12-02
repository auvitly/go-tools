package log

import (
	"context"
)

type Level int

type Logger interface {
	Log(ctx context.Context, lvl Level, msg string, args ...any)
	Fatal(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Debug(ctx context.Context, msg string, args ...any)
	With(args ...any) Logger
}

package recovery

import (
	"context"
	"github.com/auvitly/go-tools/stderrs"
)

var (
	_builder Builder
	_message = "internal server error: unhandled exception"
)

// SetMessage - set message for standard error.
func SetMessage(message string) Builder {
	return _builder.SetMessage(message)
}

// OnError - perform error enrichment.
func OnError(err *error) Builder {
	return _builder.OnError(err)
}

// WithoutHandlers - allows you to reset all handlers for the selected call.
func WithoutHandlers() Builder {
	return _builder.WithoutHandlers()
}

// On - perform standard error enrichment.
func On(err **stderrs.Error) Builder {
	return _builder.On(err)
}

// WithHandlers - add sync exception handler.
func WithHandlers(handlers ...Handler) Builder {
	return _builder.WithHandlers(handlers...)
}

// WithAsyncHandlers - add async exception handler.
func WithAsyncHandlers(handlers ...AsyncHandler) Builder {
	return _builder.WithAsyncHandlers(handlers...)
}

// Do - perform panic processing with context. Called exclusively via defer.
func Do() {
	if msg := recover(); msg != nil {
		_builder.recovery(context.Background(), msg)
	}
}

// DoContext - perform panic processing with context. Called exclusively via defer.
func DoContext(ctx context.Context) {
	if msg := recover(); msg != nil {
		_builder.recovery(ctx, msg)
	}
}

// RegistryHandlers - add handlers for global execution.
func RegistryHandlers(handlers ...Handler) {
	_builder = _builder.WithHandlers(handlers...)
}

// RegistryAsyncHandlers - add handlers for global async execution.
func RegistryAsyncHandlers(handlers ...AsyncHandler) {
	_builder = _builder.WithAsyncHandlers(handlers...)
}

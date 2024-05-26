package recovery

import (
	"context"
	"github.com/auvitly/go-tools/stderrs"
)

var (
	_builder       Builder
	_syncHandlers  []Handler
	_asyncHandlers []Handler
	_message       = "internal server error: unhandled exception"
)

// SetMessage - set message for standard error.
func SetMessage(message string) Builder { return _builder.SetMessage(message) }

// OnError - perform error enrichment.
func OnError(err *error) Builder { return _builder.OnError(err) }

// On - perform standard error enrichment.
func On(err **stderrs.Error) Builder { return _builder.On(err) }

// WithHandlers - add exception handler.
func WithHandlers(handlers ...Handler) Builder { return _builder.WithHandlers(handlers...) }

// WithAsyncHandlers - add async exception handler.
func WithAsyncHandlers(handlers ...Handler) Builder { return _builder.WithAsyncHandlers(handlers...) }

// Do - perform panic processing with context. Called exclusively via defer.
func Do(ctx context.Context) {
	if msg := recover(); msg != nil {
		_builder.recovery(ctx, msg)
	}
}

// RegistryHandlers - add handlers for global execution.
func RegistryHandlers(handlers ...Handler) {
	_syncHandlers = append(_syncHandlers, handlers...)
}

// RegistryAsyncHandlers - add handlers for global async execution.
func RegistryAsyncHandlers(handlers ...Handler) {
	_asyncHandlers = append(_asyncHandlers, handlers...)
}

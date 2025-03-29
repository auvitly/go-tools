package recovery

import (
	"github.com/auvitly/go-tools/stderrs"
)

var (
	_builder = Builder{message: _message}
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

// WithField - add a field to the target error.
func WithField(key string, value any) Builder {
	return _builder.WithField(key, value)
}

// WithFieldIf - Add a field to the target error if the condition is met.
func WithFieldIf(condition bool, key string, value any) Builder {
	return _builder.WithFieldIf(condition, key, value)
}

// WithFields - add a fields to the target error.
func WithFields(fields map[string]any) Builder {
	return _builder.WithFields(fields)
}

// WithFieldsIf - Add a fields to the target error if the condition is met.
func WithFieldsIf(condition bool, fields map[string]any) Builder {
	return _builder.WithFieldsIf(condition, fields)
}

// WithHandlers - add exception handler.
func WithHandlers(handlers ...Handler) Builder {
	return _builder.WithHandlers(handlers...)
}

// WithHandlersIf - add exception handler if the condition is met.
func WithHandlersIf(condition bool, handlers ...Handler) Builder {
	return _builder.WithHandlersIf(condition, handlers...)
}

// WithAsyncHandlers - add async exception handler.
func WithAsyncHandlers(handlers ...AsyncHandler) Builder {
	return _builder.WithAsyncHandlers(handlers...)
}

// WithAsyncHandlersIf - add async exception handler if the condition is met.
func WithAsyncHandlersIf(condition bool, handlers ...AsyncHandler) Builder {
	return _builder.WithAsyncHandlersIf(condition, handlers...)
}

// Do - perform panic processing with context. Called exclusively via defer.
func Do() {
	if msg := recover(); msg != nil {
		_builder.recovery(msg)
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

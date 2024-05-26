package recovery

import (
	"context"
	"github.com/auvitly/go-tools/stderrs"
)

var builder Builder

// SetMessage - set message for standard error.
func SetMessage(message string) Builder { return builder.SetMessage(message) }

// OnError - perform error enrichment.
func OnError(err *error) Builder { return builder.OnError(err) }

// On - perform standard error enrichment.
func On(err **stderrs.Error) Builder { return builder.On(err) }

// WithHandlers - add exception handler.
func WithHandlers(handlers ...Handler) Builder { return builder.WithHandlers(handlers...) }

func Do() { builder.Do() }

func DoContext(ctx context.Context) { builder.DoContext(ctx) }

func RegistryHandlers(handlers ...Handler) {
	_handlers = append(_handlers, handlers...)
}

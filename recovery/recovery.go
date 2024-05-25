package recovery

import "github.com/auvitly/go-tools/stderrs"

var builder Builder

func SetMessage(message string) Builder { return builder.SetMessage(message) }

func OnError(err *error) Builder { return builder.OnError(err) }

func OnStandardError(err **stderrs.Error) Builder { return builder.OnStandardError(err) }

func WithHandler(handlers ...Handler) Builder { return builder.WithHandler(handlers...) }

func Do() { builder.Do() }

func RegistryHandler(handlers ...Handler) {
	handlers = append(handlers, handlers...)
}

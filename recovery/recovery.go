package recovery

import "github.com/auvitly/go-tools/stderrs"

var builder Builder

func SetMessage(message string) Builder { return builder.SetMessage(message) }

func OnError(err *error) Builder { return builder.OnError(err) }

func On(err **stderrs.Error) Builder { return builder.On(err) }

func WithHandlers(handlers ...Handler) Builder { return builder.WithHandlers(handlers...) }

func Do() { builder.Do() }

func RegistryHandlers(handlers ...Handler) {
	handlers = append(handlers, handlers...)
}

package main

import (
	"errors"
	"github.com/auvitly/go-tools/recovery"
	"github.com/auvitly/go-tools/stderrs"
	"log/slog"
)

func syncHandler(any) error {
	return errors.New("syncHandler error: I'm the error")
}

func fn() (err *stderrs.Error) {
	defer recovery.WithHandlers(syncHandler).On(&err).Do()

	panic("I'm the exception")
}

func main() {
	err := fn()
	if err != nil {
		slog.Error(err.Error())
	}
}

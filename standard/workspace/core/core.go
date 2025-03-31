package core

import (
	"cmp"
	"context"

	"github.com/auvitly/go-tools/stderrs"
)

type Core[T, M cmp.Ordered] struct {
	dependencies Dependencies[T, M]
	config       Config
}

func New[T, M cmp.Ordered](
	dependencies Dependencies[T, M],
	config Config,
) (
	*Core[T, M], *stderrs.Error,
) {
	return &Core[T, M]{
		dependencies: dependencies,
		config:       config,
	}, nil
}

func (c *Core[T, M]) Start(ctx context.Context) *stderrs.Error {
	if c.config.PullingInterval == 0 {
		return stderrs.FailedPrecondition.SetMessage("pulling interval is 0")
	}

	return nil
}

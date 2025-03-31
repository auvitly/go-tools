package core

import (
	"cmp"
	"context"

	"github.com/auvitly/go-tools/stderrs"
)

type Core[T, M, S cmp.Ordered] struct {
	dependencies Dependencies[T, M, S]
	config       Config
}

func New[T, M, S cmp.Ordered](
	dependencies Dependencies[T, M, S],
	config Config,
) (
	*Core[T, M, S], *stderrs.Error,
) {
	return &Core[T, M, S]{
		dependencies: dependencies,
		config:       config,
	}, nil
}

func (c *Core[T, M, S]) Start(ctx context.Context) *stderrs.Error {
	if c.config.PullingInterval == 0 {
		return stderrs.FailedPrecondition.SetMessage("pulling interval is 0")
	}

	return nil
}

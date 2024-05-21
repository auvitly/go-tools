package stderrs_test

import (
	"errors"
	"github.com/auvitly/go-tools/stderrs"
	"github.com/stretchr/testify/assert"

	"io/fs"
	"testing"
)

func TestError_Is(t *testing.T) {
	t.Parallel()

	var std = stderrs.Internal.EmbedErrors(fs.ErrClosed, fs.ErrExist)

	assert.True(t, errors.Is(std, stderrs.Internal))
	assert.True(t, errors.Is(std, fs.ErrExist))
	assert.True(t, errors.Is(std, stderrs.Internal.EmbedErrors(fs.ErrClosed)))
	assert.True(t, errors.Is(std, stderrs.Internal.EmbedErrors(fs.ErrExist)))
	assert.True(t, errors.Is(std, stderrs.Internal.EmbedErrors(fs.ErrClosed, fs.ErrExist)))
	assert.False(t, errors.Is(std, stderrs.Internal.EmbedErrors(fs.ErrClosed, fs.ErrExist, fs.ErrNotExist)))
	assert.False(t, errors.Is(std, stderrs.Panic))
}

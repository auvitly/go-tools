package stderrs_test

import (
	"encoding/json"
	"github.com/auvitly/go-tools/stderrs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"io/fs"
	"testing"
)

func TestError_Is(t *testing.T) {
	t.Parallel()

	var std = stderrs.Internal.
		EmbedErrors(fs.ErrClosed, fs.ErrExist).
		EmbedErrors(stderrs.Unavailable.EmbedErrors(fs.ErrNotExist, stderrs.Aborted))

	assert.True(t, std.Is(stderrs.Internal))
	assert.True(t, std.Is(fs.ErrExist))
	assert.True(t, std.Is(stderrs.Internal.EmbedErrors(fs.ErrClosed)))
	assert.True(t, std.Is(stderrs.Internal.EmbedErrors(fs.ErrExist)))
	assert.True(t, std.Is(stderrs.Internal.EmbedErrors(fs.ErrClosed, fs.ErrExist)))
	assert.True(t, std.Is(stderrs.Internal.EmbedErrors(fs.ErrClosed, fs.ErrExist, fs.ErrNotExist, stderrs.Aborted)))
	assert.True(t, std.Is(stderrs.Unavailable.EmbedErrors(fs.ErrNotExist, stderrs.Aborted)))
	assert.True(t, std.Is(stderrs.Unavailable.EmbedErrors(fs.ErrNotExist)))

	assert.False(t, std.Is(fs.ErrPermission))
	assert.False(t, std.Is(stderrs.Panic))
	assert.False(t, std.Is(stderrs.Aborted.EmbedErrors(fs.ErrNotExist)))
	assert.False(t, std.Is(stderrs.Unavailable.EmbedErrors(fs.ErrPermission)))
}

func TestError_JSON(t *testing.T) {
	t.Parallel()

	var in = stderrs.Internal.
		EmbedErrors(fs.ErrClosed, fs.ErrExist).
		EmbedErrors(stderrs.Unavailable.EmbedErrors(fs.ErrNotExist))

	data, err := in.MarshalJSON()
	require.NoError(t, err)

	var out *stderrs.Error

	err = json.Unmarshal(data, &out)
	require.NoError(t, err)

	assert.True(t, out.Is(stderrs.Internal))
	assert.False(t, out.Is(stderrs.Undefined))

	assert.True(t, out.Contains(fs.ErrClosed))
	assert.True(t, out.Contains(stderrs.Internal))
	assert.True(t, out.Contains(stderrs.Internal.EmbedErrors(fs.ErrClosed)))
	assert.True(t, out.Contains(stderrs.Internal.EmbedErrors(fs.ErrClosed, fs.ErrExist)))
	assert.True(t, out.Contains(stderrs.Unavailable.EmbedErrors(fs.ErrNotExist)))
	assert.False(t, out.Contains(stderrs.Internal.EmbedErrors(fs.ErrClosed, fs.ErrExist, fs.ErrNotExist)))
	assert.False(t, out.Is(fs.ErrClosed))
}

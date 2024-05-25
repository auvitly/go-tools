package stderrs_test

import (
	"encoding/json"
	"github.com/auvitly/go-tools/stderrs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
	assert.True(t, std.Is(stderrs.Internal.EmbedErrors(stderrs.Unavailable, stderrs.Aborted)))
	assert.True(t, std.Is(stderrs.Internal.EmbedErrors(fs.ErrClosed, fs.ErrExist, fs.ErrNotExist, stderrs.Aborted)))
	assert.True(t, std.Is(stderrs.Unavailable.EmbedErrors(fs.ErrNotExist, stderrs.Aborted)))
	assert.True(t, std.Is(stderrs.Unavailable.EmbedErrors(fs.ErrNotExist)))

	assert.False(t, std.Is(fs.ErrPermission))
	assert.False(t, std.Is(stderrs.Panic))
	assert.False(t, std.Is(stderrs.Aborted.EmbedErrors(fs.ErrNotExist)))
	assert.False(t, std.Is(stderrs.Unavailable.EmbedErrors(fs.ErrExist)))
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
	assert.True(t, out.Is(stderrs.Unavailable))
	assert.False(t, out.Is(stderrs.Undefined))
}

func TestFrom(t *testing.T) {
	var err = status.Error(codes.Internal, "message")

	std, ok := stderrs.From(err)
	require.True(t, ok)

	require.True(t, std.Is(stderrs.Internal))
}

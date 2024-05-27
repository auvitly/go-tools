package recovery_test

import (
	"context"
	"github.com/auvitly/go-tools/recovery"
	"github.com/auvitly/go-tools/stderrs"
	"github.com/stretchr/testify/require"
	"io/fs"
	"testing"
)

func TestDo(t *testing.T) {
	t.Parallel()

	require.NotPanics(t, func() {
		defer recovery.Do(context.Background())

		panic("panic: message")
	})
}

func TestOnError(t *testing.T) {
	t.Parallel()

	var err = fs.ErrExist

	func() {
		defer recovery.OnError(&err).Do(context.Background())

		panic("panic: message")
	}()

	std, ok := stderrs.From(err)
	require.True(t, ok)
	require.True(t, std.Is(stderrs.Panic))
	require.True(t, std.Is(fs.ErrExist))
}

func TestOn(t *testing.T) {
	t.Parallel()

	var err = stderrs.Internal.EmbedErrors(fs.ErrExist).SetMessage("my message")

	func() {
		defer recovery.On(&err).Do(context.Background())

		panic("panic: message")
	}()

	std, ok := stderrs.From(err)
	require.True(t, ok)
	require.True(t, std.Is(fs.ErrExist))
	require.True(t, std.Is(stderrs.Panic))
	require.True(t, std.Is(stderrs.Internal))
	require.True(t, std.Is(stderrs.Internal.EmbedErrors(fs.ErrExist)))
}

func TestHandler(t *testing.T) {
	t.Parallel()

	var (
		actual string
		_panic = "panic"
	)

	func() {
		defer recovery.WithHandlers(
			func(_ context.Context, msg any) {
				actual = msg.(string)
			},
		).Do(context.Background())

		panic(_panic)
	}()

	require.Equal(t, _panic, actual)
}

func TestPanicInHandler(t *testing.T) {
	t.Parallel()

	var (
		err      *stderrs.Error
		_panic   = "panic"
		_message = "message"
	)

	func() {
		defer recovery.
			WithHandlers(func(_ context.Context, msg any) {
				panic(_panic)
			}).
			SetMessage(_message).
			On(&err).
			Do(context.Background())

		panic("")
	}()

	std, ok := stderrs.From(err.Embed)
	require.True(t, ok)
	require.True(t, std.Is(stderrs.Panic))
	require.Equal(t, _message, err.Message)
	require.Equal(t, _panic, std.Fields["panic"])
}

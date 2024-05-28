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
		defer recovery.OnError(&err).Do()

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
		defer recovery.On(&err).Do()

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
		defer recovery.WithSyncHandlers(
			func(msg any) error {
				actual = msg.(string)

				return nil
			},
		).Do()

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
			WithSyncHandlers(func(msg any) error {
				panic(_panic)

				return nil
			}).
			SetMessage(_message).
			On(&err).
			Do()

		panic("")
	}()

	std, ok := stderrs.From(err.Embed)
	require.True(t, ok)
	require.True(t, std.Is(stderrs.Panic))
	require.Equal(t, _message, err.Message)
	require.Equal(t, _panic, std.Fields["panic"])
}

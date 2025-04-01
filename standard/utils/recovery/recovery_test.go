package recovery_test

import (
	"fmt"
	"io/fs"
	"testing"

	"github.com/auvitly/go-tools/collection/stderrs"
	"github.com/auvitly/go-tools/utils/recovery"
	"github.com/stretchr/testify/require"
)

func TestDo(t *testing.T) {
	t.Parallel()

	require.NotPanics(t, func() {
		defer recovery.Do()

		panic("panic: message")
	})
}

func TestOnError(t *testing.T) {
	t.Parallel()

	var err = fs.ErrExist

	func() {
		defer recovery.WithField("key", "value").OnError(&err).Do()

		panic("panic: message")
	}()

	std, ok := stderrs.From(err)
	require.True(t, ok)
	require.True(t, std.Is(stderrs.Panic))
	require.True(t, std.Is(fs.ErrExist))
	require.Equal(t, std.Fields["key"], "value")
}

func TestOn(t *testing.T) {
	t.Parallel()

	var err = stderrs.Internal.
		WithField("key", "value").
		EmbedErrors(fs.ErrExist).
		SetMessage("my message")

	func() {
		defer recovery.WithField("key", "replaced").On(&err).Do()

		panic("panic: message")
	}()

	std, ok := stderrs.From(err)
	require.True(t, ok)
	require.True(t, std.Is(fs.ErrExist))
	require.True(t, std.Is(stderrs.Panic))
	require.True(t, std.Is(stderrs.Internal))
	require.True(t, std.Is(stderrs.Internal.EmbedErrors(fs.ErrExist)))
	require.Equal(t, std.Fields["key"], "replaced")
}

func TestHandler(t *testing.T) {
	t.Parallel()

	var (
		actual string
		_panic = "panic"
	)

	func() {
		defer recovery.WithHandlers(
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
		defer recovery.WithHandlers(func(msg any) error {
			panic(_panic)
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

func BenchmarkDo(b *testing.B) {
	var use = func() {
		defer recovery.Do()

		panic("")
	}

	for i := 0; i < b.N; i++ {
		use()
	}
}

func BenchmarkDefaultPanicHandler(b *testing.B) {
	var use = func() {
		defer func() {
			if msg := recover(); msg != nil {
				return
			}
		}()

		panic("")
	}

	for i := 0; i < b.N; i++ {
		use()
	}
}

func BenchmarkDoWithHandler(b *testing.B) {
	var fn = func(msg any) error {
		return fmt.Errorf("%v", msg)
	}

	var fns []recovery.Handler

	for j := 0; j < 10; j++ {
		fns = append(fns, fn)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var err error

		func() {
			defer recovery.OnError(&err).WithHandlers(fns...).Do()

			panic("")
		}()
	}
}

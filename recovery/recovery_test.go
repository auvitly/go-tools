package recovery_test

import (
	"github.com/auvitly/go-tools/recovery"
	"github.com/auvitly/go-tools/stderrs"
	"github.com/stretchr/testify/require"
	"io/fs"
	"testing"
)

func TestPanicOnError(t *testing.T) {
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

func TestPanicOnStdError(t *testing.T) {
	t.Parallel()

	var err error = stderrs.Internal.SetMessage("my message")

	func() {
		defer recovery.OnError(&err).Do()

		panic("panic: message")
	}()

	std, ok := stderrs.From(err)
	require.True(t, ok)
	require.True(t, std.Is(stderrs.Panic))
	require.True(t, std.Is(stderrs.Internal))
}

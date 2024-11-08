package labkit

import (
	"github.com/stretchr/testify/require"
	"io"
	"testing"
)

// IOReadAll - io.ReadAll without error.
func IOReadAll(t *testing.T, r io.Reader) []byte {
	t.Helper()

	data, err := io.ReadAll(r)
	require.NoError(t, err)

	return data
}

// IOCopy - io.Copy without error.
func IOCopy(t *testing.T, dst io.Writer, src io.Reader) int64 {
	t.Helper()

	n, err := io.Copy(dst, src)
	require.NoError(t, err)

	return n
}

// IOCopyN - io.CopyN without error.
func IOCopyN(t *testing.T, dst io.Writer, src io.Reader, n int64) int64 {
	t.Helper()

	written, err := io.CopyN(dst, src, n)
	require.NoError(t, err)

	return written
}

// IOCopyBuffer - io.CopyBuffer without error.
func IOCopyBuffer(t *testing.T, dst io.Writer, src io.Reader, buf []byte) int64 {
	t.Helper()

	written, err := io.CopyBuffer(dst, src, buf)
	require.NoError(t, err)

	return written
}

// IOReadFull - io.ReadFull without error.
func IOReadFull(t *testing.T, r io.Reader, buf []byte) int {
	t.Helper()

	written, err := io.ReadFull(r, buf)
	require.NoError(t, err)

	return written
}

// IOReadAtLeast - io.ReadAtLeast without error.
func IOReadAtLeast(t *testing.T, r io.Reader, buf []byte, min int) int {
	t.Helper()

	written, err := io.ReadAtLeast(r, buf, min)
	require.NoError(t, err)

	return written
}

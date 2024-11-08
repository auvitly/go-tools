package labkit

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func OSOpen(t *testing.T, name string) *os.File {
	t.Helper()

	file, err := os.Open(name)
	require.NoError(t, err)

	return file
}

func OSOpenFile(t *testing.T, name string, flag int, perm os.FileMode) *os.File {
	t.Helper()

	file, err := os.OpenFile(name, flag, perm)
	require.NoError(t, err)

	return file
}

func OSRemove(t *testing.T, name string) {
	t.Helper()

	err := os.Remove(name)
	require.NoError(t, err)
}

func OSCreate(t *testing.T, name string) *os.File {
	t.Helper()

	file, err := os.Create(name)
	require.NoError(t, err)

	return file
}

func OSChdir(t *testing.T, name string) {
	t.Helper()

	err := os.Chdir(name)
	require.NoError(t, err)
}

func OSChmod(t *testing.T, name string, mode os.FileMode) {
	t.Helper()

	err := os.Chmod(name, mode)
	require.NoError(t, err)
}

func OSGetwd(t *testing.T) string {
	t.Helper()

	dir, err := os.Getwd()
	require.NoError(t, err)

	return dir
}

func OSGetgroups(t *testing.T) []int {
	t.Helper()

	list, err := os.Getgroups()
	require.NoError(t, err)

	return list
}

func OSFindProcess(t *testing.T, pid int) *os.Process {
	t.Helper()

	process, err := os.FindProcess(pid)
	require.NoError(t, err)

	return process
}

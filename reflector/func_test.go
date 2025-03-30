package reflector_test

import (
	"testing"

	"github.com/auvitly/go-tools/reflector"
	"github.com/stretchr/testify/require"
)

func TestParseFunc(t *testing.T) {
	fn, err := reflector.ParseFunc(TestParseFunc)
	require.NoError(t, err)

	t.Log(fn)
}

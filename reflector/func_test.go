package reflector_test

import (
	"github.com/auvitly/go-tools/reflector"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestScanFunction(t *testing.T) {
	var A = func() {}

	fn, err := reflector.ScanFunc(A)
	require.NoError(t, err)

	t.Log(fn.Runtime.Name())
}

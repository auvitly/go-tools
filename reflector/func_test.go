package reflector_test

import (
	"testing"

	"github.com/auvitly/go-tools/reflector"
	"github.com/stretchr/testify/require"
)

func TestScanFunction(t *testing.T) {
	fn, err := reflector.ScanFunc(TestScanFunction)
	require.NoError(t, err)

	t.Log(fn)
}

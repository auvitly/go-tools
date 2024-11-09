package lab_test

import (
	"github.com/auvitly/go-tools/lab"
	"github.com/auvitly/go-tools/lab/vault"
	"github.com/stretchr/testify/require"
	"net"
	"testing"
)

func TestStoreLoad(t *testing.T) {
	t.Parallel()

	var tests = []lab.Test[
		net.HardwareAddr,
		net.HardwareAddr,
	]{
		{
			Name: "#1:Storage:EqualValues",
			In:   vault.Store(t, "mac", lab.Value(net.ParseMAC("b1:b2:1e:68:ab:d4"))),
			Out:  vault.Load[net.HardwareAddr](t, "mac"),
		},
	}

	for i := range tests {
		var test = tests[i]

		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, test.In, test.Out)
		})
	}
}

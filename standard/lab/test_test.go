package lab_test

import (
	"net"
	"testing"
	"time"

	"github.com/auvitly/go-tools/lab"
	"github.com/auvitly/go-tools/lab/vault"
	"github.com/stretchr/testify/require"
)

func TestVault(t *testing.T) {
	t.Parallel()

	var tests = []lab.Test[
		lab.Any,
		lab.Any,
	]{
		{
			Name: "net.HardwareAddr",
			In:   vault.Store(t, "mac", lab.Glue(net.ParseMAC("b1:b2:1e:68:ab:d4"))[0]),
			Out:  vault.Load[net.HardwareAddr](t, "mac"),
		},
		{
			Name: "net.IP",
			In:   vault.Store(t, "ip", lab.Glue(net.ParseCIDR("127.0.0.1/24"))[0]),
			Out:  vault.Load[net.IP](t, "ip"),
		},
		{
			Name: "*net.IPNet",
			In:   vault.Store(t, "cidr", lab.Glue(net.ParseCIDR("127.0.0.1/24"))[1]),
			Out:  vault.Load[*net.IPNet](t, "cidr"),
		},
		{
			Name: "time.Time",
			In:   vault.Store(t, "time", lab.First(time.Parse(time.RFC3339, "2006-01-02T15:04:05+07:00"))),
			Out:  vault.Load[time.Time](t, "time"),
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

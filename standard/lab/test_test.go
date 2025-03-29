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

	var v = vault.New()

	var tests = []lab.Test[
		lab.Any,
		lab.Any,
	]{
		{
			Name: "net.HardwareAddr",
			In:   vault.Store(v, "mac", lab.Return[net.HardwareAddr](net.ParseMAC("b1:b2:1e:68:ab:d4"))(lab.FirstValue)),
			Out:  vault.Load[net.HardwareAddr](v, "mac"),
		},
		{
			Name: "net.IP",
			In:   vault.Store(v, "ip", lab.Return[net.IP](net.ParseCIDR("127.0.0.1/24"))(lab.FirstValue)),
			Out:  vault.Load[net.IP](v, "ip"),
		},
		{
			Name: "*net.IPNet",
			In:   vault.Store(v, "cidr", lab.Return[*net.IPNet](net.ParseCIDR("127.0.0.1/24"))(lab.SecondValue)),
			Out:  vault.Load[*net.IPNet](v, "cidr"),
		},
		{
			Name: "time.Time",
			In:   vault.Store(v, "time", lab.First(time.Parse(time.RFC3339, "2006-01-02T15:04:05+07:00"))),
			Out:  vault.Load[time.Time](v, "time"),
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

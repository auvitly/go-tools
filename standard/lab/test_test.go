package lab_test

import (
	"github.com/auvitly/go-tools/lab"
	"github.com/auvitly/go-tools/lab/kit"
	"github.com/auvitly/go-tools/lab/vault"
	"github.com/stretchr/testify/require"
	"io"
	"net"
	"testing"
	"time"
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
			In:   vault.Store[net.HardwareAddr](v, "mac", lab.Value(net.ParseMAC("b1:b2:1e:68:ab:d4"))),
			Out:  vault.Load[net.HardwareAddr](v, "mac"),
		},
		{
			Name: "net.IP",
			In:   vault.Store[net.IP](v, "ip", lab.FirstValue(net.ParseCIDR("127.0.0.1/24"))),
			Out:  vault.Load[net.IP](v, "ip"),
		},
		{
			Name: "*net.IPNet",
			In:   vault.Store[*net.IPNet](v, "cidr", lab.SecondValue(net.ParseCIDR("127.0.0.1/24"))),
			Out:  vault.Load[*net.IPNet](v, "cidr"),
		},
		{
			Name: "time.Time",
			In:   vault.Store[time.Time](v, "time", lab.Value(time.Parse(time.RFC3339, "2006-01-02T15:04:05+07:00"))),
			Out:  vault.Load[time.Time](v, "time"),
		},
		{
			Name: "kit.Now",
			In:   vault.Store[time.Time](v, "now", kit.Now),
			Out:  vault.Load[time.Time](v, "now"),
		},
		{
			Name: "error",
			In:   vault.Store[error](v, "error", io.ErrClosedPipe),
			Out:  vault.Load[error](v, "error"),
		},
		{
			Name: "kit.IPv4",
			In:   vault.Store[net.IP](v, "kit.IPv4", kit.IPv4),
			Out:  vault.Load[net.IP](v, "kit.IPv4"),
		},
		{
			Name: "kit.IPv6",
			In:   vault.Store[net.IP](v, "kit.IPv6", kit.IPv6),
			Out:  vault.Load[net.IP](v, "kit.IPv6"),
		},
		{
			Name: "kit.Int",
			In:   vault.Store[int](v, "kit.Int", kit.Int),
			Out:  vault.Load[int](v, "kit.Int"),
		},
		{
			Name: "kit.Int8",
			In:   vault.Store[int8](v, "kit.Int8", kit.Int8),
			Out:  vault.Load[int8](v, "kit.Int8"),
		},
		{
			Name: "kit.String",
			In:   vault.Store[string](v, "kit.String", kit.String),
			Out:  vault.Load[string](v, "kit.String"),
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

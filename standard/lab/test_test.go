package lab_test

import (
	"github.com/auvitly/go-tools/lab"
	"github.com/auvitly/go-tools/lab/vault"
	"github.com/stretchr/testify/require"
	"io"
	"net"
	"testing"
	"time"
)

func TestVault(t *testing.T) {
	t.Parallel()

	var tests = []lab.Test[
		lab.TODO,
		lab.TODO,
	]{
		{
			Name: "#1:Vault:EqualValues:net.HardwareAddr",
			In:   vault.Store(t, "mac", lab.Value(net.ParseMAC("b1:b2:1e:68:ab:d4"))),
			Out:  vault.Load[net.HardwareAddr](t, "mac"),
		},
		{
			Name: "#2:Vault:EqualValues:net.IP",
			In:   vault.Store(t, "ip", lab.FirstValue(net.ParseCIDR("127.0.0.1/24"))),
			Out:  vault.Load[net.IP](t, "ip"),
		},
		{
			Name: "#3:Vault:EqualValues:*net.IPNet",
			In:   vault.Store(t, "cidr", lab.SecondValue(net.ParseCIDR("127.0.0.1/24"))),
			Out:  vault.Load[*net.IPNet](t, "cidr"),
		},
		{
			Name: "#4:Vault:EqualValues:time.Time",
			In:   vault.Store(t, "time", lab.Value(time.Parse(time.RFC3339, "2006-01-02T15:04:05+07:00"))),
			Out:  vault.Load[time.Time](t, "time"),
		},
		{
			Name: "#5:Vault:EqualValues:lab.Now",
			In:   vault.Store(t, "now", lab.Now),
			Out:  vault.Load[time.Time](t, "now"),
		},
		{
			Name: "#6:Vault:EqualValues:error",
			In:   vault.Store(t, "error", io.ErrClosedPipe),
			Out:  vault.Load[error](t, "error"),
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

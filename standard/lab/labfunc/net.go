package labfunc

import (
	"github.com/stretchr/testify/require"
	"net"
	"testing"
)

// NetParseIP - net.ParseIP without error. Returns pointer.
func NetParseIP(t *testing.T, s string) *net.IP {
	t.Helper()

	ip := net.ParseIP(s)
	require.NotNil(t, ip)

	return &ip
}

// NetParseMAC - net.ParseMAC without error. Returns pointer.
func NetParseMAC(t *testing.T, s string) *net.HardwareAddr {
	t.Helper()

	mac, err := net.ParseMAC(s)
	require.NoError(t, err)

	return &mac
}

// NetParseCIDR - net.ParseCIDR without error.
func NetParseCIDR(t *testing.T, s string) *net.IPNet {
	t.Helper()

	_, cidr, err := net.ParseCIDR(s)
	require.NoError(t, err)
	require.NotNil(t, cidr)

	return cidr
}

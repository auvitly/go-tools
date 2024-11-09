package random

import (
	"net"
)

var v4InV6Prefix = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff}

func IPv4() (ip net.IP) {
	return append(v4InV6Prefix, Bytes(4)...)
}

func IPv6() (ip net.IP) {
	return Bytes(16)
}

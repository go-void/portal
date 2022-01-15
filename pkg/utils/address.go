package utils

import (
	"net"
)

// DNSAddress returns an address with the following format: <ip-address>:53
func DNSAddress(ip net.IP) string {
	return net.JoinHostPort(ip.String(), "53")
}

package utils

import (
	"fmt"
	"net"
)

// DNSAddress returns an address with the following format: <ip-address>:53
func DNSAddress(ip net.IP) string {
	if len(ip) == 16 {
		return fmt.Sprintf("%s:%d", ip.String(), 53)
	} else {
		return fmt.Sprintf("[%s]:%d", ip.String(), 53)
	}
}

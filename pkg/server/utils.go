package server

import (
	"net"
	"net/netip"
)

// See https://github.com/golang/go/issues/49097 for upcoming Dialer changes

// createByteBuffer returns a function which creates a slice of bytes with the provided length
func createByteBuffer(size int) func() interface{} {
	return func() interface{} {
		return make([]byte, size)
	}
}

// createUDPListener creates a UDP listener
func createUDPListener(network string, addrPort netip.AddrPort) (*net.UDPConn, error) {
	return net.ListenUDP(network, net.UDPAddrFromAddrPort(addrPort))
}

// createTCPListener creates a TCP listener
func createTCPListener(network string, addrPort netip.AddrPort) (*net.TCPListener, error) {
	return net.ListenTCP(network, net.TCPAddrFromAddrPort(addrPort))
}

package server

import "net"

// createByteBuffer returns a function which creates a slice of bytes with the provided length
func createByteBuffer(size int) func() interface{} {
	return func() interface{} {
		return make([]byte, size)
	}
}

// createUDPListener creates a UDP listener
func createUDPListener(network string, address net.IP, port int) (*net.UDPConn, error) {
	return net.ListenUDP(network, &net.UDPAddr{
		IP:   address,
		Port: port,
	})
}

// createTCPListener creates a TCP listener
func createTCPListener(network string, address net.IP, port int) (*net.TCPListener, error) {
	return net.ListenTCP(network, &net.TCPAddr{
		IP:   address,
		Port: port,
	})
}

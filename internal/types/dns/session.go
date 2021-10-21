package dns

import "net"

// Session holds the remote IP address, remote port
// and optional additional data of a TCP or UDP
// "connection"
type Session struct {
	Address    *net.UDPAddr
	Additional []byte
}

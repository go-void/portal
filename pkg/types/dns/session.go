package dns

import (
	"net/netip"
)

// Session holds the remote IP address, remote port and optional additional data of a UDP "connection"
type Session struct {
	AddrPort   netip.AddrPort
	Additional []byte
}

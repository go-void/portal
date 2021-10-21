package reader

import (
	"net"

	"github.com/go-void/portal/internal/types/dns"
)

// Reader defines a reader interface which enbables implementers
// to read incoming TCP and UDP messages from listeners and return
// these messages as a byte slice, session data like remote IP
// address and additional data or an error.
type Reader interface {
	// ReadUDP reads a single message from a UDP connection and
	// returns the message as a slice of bytes, the length of the
	// read message, session data or an error
	ReadUDP(*net.UDPConn, []byte) (int, dns.Session, error)
}

type DefaultReader struct {
	OOBSize int
}

func NewDefault(oobSize int) Reader {
	return &DefaultReader{
		OOBSize: oobSize,
	}
}

// ReadUDP reads a single message from a UDP connection and
// returns the message as a slice of bytes, the length of the
// read message, session data or an error
func (r *DefaultReader) ReadUDP(c *net.UDPConn, message []byte) (int, dns.Session, error) {
	oob := make([]byte, r.OOBSize)
	session := dns.Session{}

	mn, oobn, _, addr, err := c.ReadMsgUDP(message, oob)
	if err != nil {
		return mn, session, err
	}

	session.Address = addr
	session.Additional = oob[:oobn]
	return mn, session, nil
}

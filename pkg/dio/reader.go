package dio

import (
	"encoding/binary"
	"io"
	"net"

	"github.com/go-void/portal/pkg/types/dns"
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

	ReadTCP(net.Conn) ([]byte, error)
}

type DefaultReader struct {
	AncillarySize int
}

func NewDefaultReader(ancillarySize int) Reader {
	return &DefaultReader{
		AncillarySize: ancillarySize,
	}
}

// ReadUDP reads a single message from a UDP connection and returns the message as a slice of bytes, the length of the
// read message, session data or an error
func (r *DefaultReader) ReadUDP(c *net.UDPConn, message []byte) (int, dns.Session, error) {
	ancillary := make([]byte, r.AncillarySize)
	session := dns.Session{}

	messageLen, ancillaryLen, _, addr, err := c.ReadMsgUDP(message, ancillary)
	if err != nil {
		return messageLen, session, err
	}

	session.Address = addr
	session.Additional = ancillary[:ancillaryLen]
	return messageLen, session, nil
}

func (r *DefaultReader) ReadTCP(conn net.Conn) ([]byte, error) {
	var length uint16

	err := binary.Read(conn, binary.BigEndian, &length)
	if err != nil {
		return nil, err
	}

	message := make([]byte, length)
	_, err = io.ReadFull(conn, message)
	if err != nil {
		return nil, err
	}

	return message, nil
}

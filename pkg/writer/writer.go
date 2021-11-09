// Package writer provides functions to write DNS messages (byte slices) back to clients
package writer

import (
	"encoding/binary"
	"net"
)

// Writer describes an interface which allows to write DNS messages (byte slices) back to a client via UDP and TCP
type Writer interface {
	// WriteUDPClose writes a byte slice back to a client with
	// 'addr' via the provided UDP conn and closes it afterwards
	WriteUDPClose(*net.UDPConn, []byte, *net.UDPAddr) error

	// WriteUDP writes a byte slice back to a client with
	// 'addr' via the provided UDP conn
	WriteUDP(*net.UDPConn, []byte, *net.UDPAddr) error

	// WriteTCPClose writes a byte slice back to a client with
	// 'addr' via the provided TCP conn and closes it afterwards
	WriteTCPClose(*net.TCPConn, []byte) error

	// WriteTCP writes a byte slice back to a client with
	// 'addr' via the provided TCP conn
	WriteTCP(*net.TCPConn, []byte) error
}

// DefaultWriter is the default implementation of the 'Writer' interface
type DefaultWriter struct {
}

// NewDefault returns a new default writer
func NewDefault() *DefaultWriter {
	return &DefaultWriter{}
}

// WriteUDPClose writes a byte slice back to a client with 'addr' via the provided UDP conn and closes it afterwards
func (w *DefaultWriter) WriteUDPClose(conn *net.UDPConn, buf []byte, addr *net.UDPAddr) error {
	err := w.WriteUDP(conn, buf, addr)
	if err != nil {
		return err
	}
	return conn.Close()
}

// WriteUDP writes a byte slice back to a client with 'addr' via the provided UDP conn
func (w *DefaultWriter) WriteUDP(conn *net.UDPConn, buf []byte, addr *net.UDPAddr) error {
	_, err := conn.WriteToUDP(buf, addr)
	if err != nil {
		return err
	}
	return nil
}

// WriteTCPClose writes a byte slice back to a client with 'addr' via the provided TCP conn and closes it afterwards
func (w *DefaultWriter) WriteTCPClose(conn *net.TCPConn, buf []byte) error {
	err := w.WriteTCP(conn, buf)
	if err != nil {
		return err
	}
	return conn.Close()
}

// WriteTCP writes a byte slice back to a client with 'addr' via the provided TCP conn
func (w *DefaultWriter) WriteTCP(conn *net.TCPConn, buf []byte) error {
	b := make([]byte, len(buf)+2)
	binary.BigEndian.PutUint16(b, uint16(len(b)))
	copy(b[2:], buf)

	_, err := conn.Write(b)
	if err != nil {
		return err
	}
	return nil
}

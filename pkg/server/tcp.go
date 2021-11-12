package server

import (
	"fmt"
	"net"

	"github.com/go-void/portal/pkg/types/dns"
)

// serveTCP is the main listen / answer loop, which handles DNS queries and responses via TCP
func (s *Server) serveTCP() error {
	for s.isRunning() {
		conn, err := s.TCPListener.AcceptTCP()
		if err != nil {
			return err
		}

		b, err := s.Reader.ReadTCP(conn)
		if err != nil {
			return err
		}

		header, offset, err := s.Unpacker.UnpackHeader(b)
		if err != nil {
			return err
		}

		switch s.AcceptFunc(header) {
		case AcceptMessage:
			m, err := s.Unpacker.Unpack(header, b, offset)
			if err != nil {
				return err
			}

			s.wg.Add(1)
			go s.handleTCP(m, conn)
		}
	}

	return nil
}

// handleTCP handles name matching and returns a response message via TCP
func (s *Server) handleTCP(message dns.Message, conn *net.TCPConn) {
	message, err := s.handle(message)
	if err != nil {
		fmt.Println(err)
		return
	}

	s.writeTCP(message, conn)
}

// writeTCP packs a DNS message and writes it back to the requesting DNS client via TCP
func (s *Server) writeTCP(message dns.Message, conn *net.TCPConn) {
	defer s.wg.Done()

	b, err := s.Packer.Pack(message)
	if err != nil {
		// Handle
		return
	}

	err = s.Writer.WriteTCPClose(conn, b)
	if err != nil {
		// Handle
		fmt.Println(err)
		return
	}
}
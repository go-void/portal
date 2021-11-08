package server

import (
	"fmt"

	"github.com/go-void/portal/pkg/types/dns"
)

// serveUDP is the main listen / answer loop, which handles DNS queries and responses via UDP
func (s *Server) serveUDP() error {
	// FIXME (Techassi): Handle shutdown with shutdown context and signals
	for s.isRunning() {
		b, session, err := s.readUDP()
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
			go s.handleUDP(m, session)
		}
	}

	return nil
}

// handleUDP handles name matching and returns a response message via UDP
func (s *Server) handleUDP(message dns.Message, session dns.Session) {
	message, err := s.handle(message)
	if err != nil {
		fmt.Println(err)
		return
	}

	s.writeUDP(message, session)
}

// readUDP reads a UDP message from the UDP connection by retrieving a byte buffer from the message pool
func (s *Server) readUDP() ([]byte, dns.Session, error) {
	rm := s.messageList.Get().([]byte)
	mn, session, err := s.Reader.ReadUDP(s.UDPListener, rm)
	if err != nil {
		s.messageList.Put(rm)
		return nil, session, err
	}
	return rm[:mn], session, nil
}

// writeUDP packs a DNS message and writes it back to the requesting DNS client via UDP
func (s *Server) writeUDP(message dns.Message, session dns.Session) {
	defer s.wg.Done()

	b, err := s.Packer.Pack(message)
	if err != nil {
		// Handle
		return
	}

	_, err = s.UDPListener.WriteToUDP(b, session.Address)
	if err != nil {
		// Handle
		return
	}
}

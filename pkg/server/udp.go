package server

import (
	"github.com/go-void/portal/pkg/logger"
	"github.com/go-void/portal/pkg/types/dns"

	"go.uber.org/zap"
)

// serveUDP is the main listen / answer loop, which handles DNS queries and
// responses via UDP
func (s *Server) serveUDP() {
	s.Logger.Info("start UDP listener", zap.String("context", "server"))

	// FIXME (Techassi): Handle shutdown with shutdown context and signals
	for s.isRunning() {
		b, session, err := s.readUDP()
		if err != nil {
			s.Logger.Error(logger.ErrUDPRead,
				zap.String("context", "server"),
				zap.Error(err),
			)
		}

		header, offset, err := s.Unpacker.UnpackHeader(b)
		if err != nil {
			s.Logger.Error(logger.ErrUnpackDNSHeader,
				zap.String("context", "server"),
				zap.Error(err),
			)
		}

		switch result := s.AcceptFunc(header); result {
		case AcceptMessage:
			m, err := s.Unpacker.Unpack(header, b, offset)
			if err != nil {
				s.Logger.Error(logger.ErrUnpackDNSMessage,
					zap.String("context", "server"),
					zap.Error(err),
				)
			}

			s.conns.Add(1)
			go s.handleUDP(m, session)
		default:
			// TODO (Techassi): Handle with appropriate response
			s.Logger.Info(logger.ErrAcceptMessage,
				zap.String("context", "server"),
				zap.String("reason", result.String()),
				zap.Error(err),
			)
			panic("Not accepted")
		}
	}
}

// handleUDP handles name matching and returns a response message via UDP
func (s *Server) handleUDP(message *dns.Message, session dns.Session) {
	message, err := s.handle(message, session.AddrPort)
	if err != nil {
		s.Logger.Error(logger.ErrHandleRequest,
			zap.String("context", "server"),
			zap.Error(err),
		)
		return
	}

	s.writeUDP(message, session)
}

// readUDP reads a UDP message from the UDP connection by retrieving a byte
// buffer from the message pool
func (s *Server) readUDP() ([]byte, dns.Session, error) {
	rm := s.messageList.Get().([]byte)
	mn, session, err := s.Reader.ReadUDP(s.UDPListener, rm)
	if err != nil {
		s.messageList.Put(&rm)
		return nil, session, err
	}
	return rm[:mn], session, nil
}

// writeUDP packs a DNS message and writes it back to the requesting DNS client
// via UDP
func (s *Server) writeUDP(message *dns.Message, session dns.Session) {
	defer s.conns.Done()

	b, err := s.Packer.Pack(message)
	if err != nil {
		s.Logger.Error(logger.ErrPackDNSMessage,
			zap.String("context", "server"),
			zap.Error(err),
		)
		return
	}

	err = s.Writer.WriteUDP(s.UDPListener, b, session.AddrPort)
	if err != nil {
		s.Logger.Error(logger.ErrUDPWrite,
			zap.String("context", "server"),
			zap.Error(err),
		)
		return
	}
}

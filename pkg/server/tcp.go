package server

import (
	"net"

	"github.com/go-void/portal/pkg/logger"
	"github.com/go-void/portal/pkg/types/dns"

	"go.uber.org/zap"
)

// serveTCP is the main listen / answer loop, which handles DNS queries and
// responses via TCP
func (s *Server) serveTCP() {
	for s.isRunning() {
		conn, err := s.TCPListener.AcceptTCP()
		if err != nil {
			s.Logger.Error(logger.ErrTCPAccept,
				zap.String("context", "server"),
				zap.Error(err),
			)
		}

		b, err := s.Reader.ReadTCP(conn)
		if err != nil {
			s.Logger.Error(logger.ErrTCPRead,
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
			go s.handleTCP(m, conn)
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

// handleTCP handles name matching and returns a response message via TCP
func (s *Server) handleTCP(message *dns.Message, conn *net.TCPConn) {
	addr := conn.RemoteAddr().(*net.TCPAddr)
	message, err := s.handle(message, addr.IP)
	if err != nil {
		s.Logger.Error(logger.ErrHandleRequest,
			zap.String("context", "server"),
			zap.Error(err),
		)
		return
	}

	s.writeTCP(message, conn)
}

// writeTCP packs a DNS message and writes it back to the requesting DNS client via TCP
func (s *Server) writeTCP(message *dns.Message, conn *net.TCPConn) {
	defer s.conns.Done()

	b, err := s.Packer.Pack(message)
	if err != nil {
		s.Logger.Error(logger.ErrPackDNSMessage,
			zap.String("context", "server"),
			zap.Error(err),
		)
		return
	}

	err = s.Writer.WriteTCPClose(conn, b)
	if err != nil {
		s.Logger.Error(logger.ErrTCPWriteClose,
			zap.String("context", "server"),
			zap.Error(err),
		)
		return
	}
}

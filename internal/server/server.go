// Package server provides functions to create and manage DNS server instances
package server

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/go-void/portal/internal/cache"
	"github.com/go-void/portal/internal/config"
	"github.com/go-void/portal/internal/constants"
	"github.com/go-void/portal/internal/filter"
	"github.com/go-void/portal/internal/pack"
	"github.com/go-void/portal/internal/reader"
	"github.com/go-void/portal/internal/resolver"
	"github.com/go-void/portal/internal/store"
	"github.com/go-void/portal/internal/types/dns"
	"github.com/go-void/portal/internal/types/rr"

	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

var (
	ErrServerAlreadyRunning = errors.New("server already running")
	ErrNoSuchNetwork        = errors.New("no such network")
	ErrNoQuestions          = errors.New("no questions")
)

// Server describes options for running a DNS server
type Server struct {
	// The address the server is running on (default: 0.0.0.0)
	Address net.IP

	// Network the server is using (default: udp)
	Network string

	// Port the server is listening on (default: 53)
	Port int

	// If the server is running in TCP mode, this listener
	// listens for incoming DNS messages
	TCPListener *net.TCPListener

	// If the server us running in UDP mode, this listener
	// listens for incoming UDP messages
	UDPListener *net.UDPConn

	// Reader implements the Reader interface to read
	// incoming TCP and UDP messages
	Reader reader.Reader

	// Unpacker implements the Unpacker interface to unwrap
	// DNS messages
	Unpacker pack.Unpacker

	// Packer implements the Packer interface to pack
	// DNS messages
	Packer pack.Packer

	// Filter implements the Filter interface to enable DNS
	// filtering
	Filter filter.Filter

	// Cache implements the Cache interface to store record
	// data in memory for a specified TTL
	Cache cache.Cache

	// Store implements the Store interface to get and set
	// record data. A store can be anything, e.g. a file or
	// a database. The store has access to a cache instance
	// which stores record data in memory. The default store
	// does not use the cache as it is implemented as a
	// in-memory tree structure
	Store store.Store

	// Resolver implements the Resolver interface to resolve
	// unknown domain names either iteratively or recursively.
	// After resolving the record data should get stored in
	// the cache
	Resolver resolver.Resolver

	// UDPMessageSize is the default message size to create
	// the temporary slice of bytes within the messageList
	// pool
	UDPMessageSize int

	// AncillarySize is the maximum size required to store
	// ancillary data of UDP messages. This is culculated
	// beforehand to avoid re-calculation of the same size
	// over and over again
	AncillarySize int

	// AcceptFunc validates if the DNS message should be
	// accepted or rejected (when sending invalid data)
	AcceptFunc AcceptFunc

	// messageList contains short-lived byte-slices in which
	// incoming messages get stored for the duration of the
	// question / answer cycle
	messageList sync.Pool

	// RW mutex to lock and unlock while writing and
	// reading server data
	lock sync.RWMutex

	wg sync.WaitGroup

	// This indicates if the server instance is running
	running bool

	// This indicates if the server / store is using a cache
	usesCache bool
}

// New creates a new DNS server instance
func New() *Server {
	server := &Server{
		UDPMessageSize: constants.UDPMinMessageSize,
		messageList:    sync.Pool{},
		wg:             sync.WaitGroup{},
	}
	server.messageList.New = createByteBuffer(constants.UDPMinMessageSize)

	return server
}

// ListenAndServe starts the listen / respond loop of DNS messages
func (s *Server) ListenAndServe() error {
	if s.isRunning() {
		return ErrServerAlreadyRunning
	}

	s.running = true

	switch s.Network {
	case "udp", "udp4", "udp6":
		listener, err := createUDPListener(s.Network, s.Address, s.Port)
		if err != nil {
			return err
		}
		s.UDPListener = listener
		return s.serveUDP()
	case "tcp", "tcp4", "tcp6":
		fmt.Println("TCP!")
		listener, err := createTCPListener(s.Network, s.Address, s.Port)
		if err != nil {
			return err
		}
		s.TCPListener = listener
		return s.serveTCP()
	}

	return ErrNoSuchNetwork
}

// Configure initializes default server parameters and checks if all neccesary components are registered. If not it
// falls back to defaults. This functions expects a validated config.Config
func (s *Server) Configure(c *config.Config) {
	ancillary4 := ipv4.NewControlMessage(ipv4.FlagDst | ipv4.FlagInterface)
	ancillary6 := ipv6.NewControlMessage(ipv6.FlagDst | ipv6.FlagInterface)

	if len(ancillary4) > len(ancillary6) {
		s.AncillarySize = len(ancillary4)
	} else {
		s.AncillarySize = len(ancillary6)
	}

	s.Address = c.Server.Address
	s.Network = c.Server.Network
	s.Port = c.Server.Port

	if s.Filter == nil {
		s.Filter = filter.New()
	}

	if s.Cache == nil {
		s.Cache = cache.NewDefaultCache()
		s.usesCache = true
	}

	if s.Resolver == nil {
		switch c.Resolver.Mode {
		case "r":
			// TODO (Techassi): Adjust API, provide path to hints file instead of the hints itself
			s.Resolver = resolver.NewRecursiveResolver([]net.IP{}, s.Cache)
		case "i":
			s.Resolver = resolver.NewIterativeResolver()
		case "f":
			s.Resolver = resolver.NewForwardingResolver(c.Resolver.Upstream, s.Cache)
		}
	}

	if s.Unpacker == nil {
		s.Unpacker = pack.NewDefaultUnpacker()
	}

	if s.Packer == nil {
		s.Packer = pack.NewDefaultPacker()
	}

	if s.Reader == nil {
		s.Reader = reader.NewDefault(s.AncillarySize)
	}

	if s.AcceptFunc == nil {
		s.AcceptFunc = DefaultAcceptFunc
	}

	if s.Store == nil {
		s.Store = store.NewDefault()
	}
}

// serveUDP is the main listen / answer loop, which handles DNS queries and responses via UDP
func (s *Server) serveUDP() error {
	// FIXME (Techassi): Handle shutdown with shutdown context and signals
	for s.isRunning() {
		b, session, err := s.readUDP()
		if err != nil {
			return err
		}

		fmt.Println(b)

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

// serveTCP is the main listen / answer loop, which handles DNS queries and responses via TCP
func (s *Server) serveTCP() error {
	for s.isRunning() {
		conn, err := s.TCPListener.Accept()
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

// handleUDP handles name matching and returns a response message via UDP
func (s *Server) handleUDP(message dns.Message, session dns.Session) {
	message, err := s.handle(message)
	if err != nil {
		fmt.Println(err)
		return
	}

	s.writeUDPMessage(message, session)
}

// handleTCP handles name matching and returns a response message via TCP
func (s *Server) handleTCP(message dns.Message, conn net.Conn) {
	message, err := s.handle(message)
	if err != nil {
		fmt.Println(err)
		return
	}

	s.writeTCPMessage(message, conn)
}

// handle handles name matching and returns a response message
func (s *Server) handle(message dns.Message) (dns.Message, error) {
	if len(message.Question) == 0 {
		s.wg.Done()
		return message, ErrNoQuestions
	}

	// TODO (Techassi): Add support for ANY queries
	var err error

	filtered, message, err := s.Filter.Match(message)
	if err != nil {
		// FIXME (Techassi): How whould we handle a filter error? Should we abort or continue (and answer the query)
		return message, err
	}

	if filtered {
		return message, nil
	}

	// TODO (Techassi): Clean this up. Is there a more elegant solution?
	var status cache.Status
	var entry cache.Entry
	var record rr.RR

	if s.usesCache {
		entry, status, err = s.Cache.LookupQuestion(message.Question[0])
		record = entry.Record
	}

	if status != cache.Hit || !s.usesCache {
		record, err = s.Store.Get(message.Question[0])
	}

	if err != nil {
		record, err = s.Resolver.ResolveQuestion(message.Question[0])
		if err != nil && !errors.Is(err, resolver.ErrNoAnswer) {
			return message, err
		}
	}

	message.AddAnswer(record)
	message.IsResponse()
	return message, nil
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

// writeUDPMessage packs a DNS message and writes it back to the requesting DNS client via UDP
func (s *Server) writeUDPMessage(message dns.Message, session dns.Session) {
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

// writeTCPMessage packs a DNS message and writes it back to the requesting DNS client via TCP
func (s *Server) writeTCPMessage(message dns.Message, conn net.Conn) {
	defer func() {
		s.wg.Done()
		// conn.Close()
	}()

	b, err := s.Packer.Pack(message)
	if err != nil {
		// Handle
		return
	}

	// TODO (Techassi): Create Writer interface which abstracts writing messages
	m := make([]byte, len(b)+2)
	binary.BigEndian.PutUint16(m, uint16(len(b)))
	copy(m[2:], b)

	_, err = conn.Write(m)
	if err != nil {
		// Handle
		fmt.Println(err)
		return
	}
}

// isRunning returns if the server instance is running
func (s *Server) isRunning() bool {
	s.lock.RLock()
	running := s.running
	s.lock.RUnlock()
	return running
}

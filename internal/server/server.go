// Package server provides functions to create and manage DNS
// server instances
package server

import (
	"errors"
	"net"
	"sync"

	"github.com/go-void/portal/internal/cache"
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
	ErrServerAlreadyRunning = errors.New("Server already running")
	ErrNoSuchNetwork        = errors.New("No such network")
)

// Server describes options for running a DNS server
type Server struct {
	// The address the server is running on (default: 0.0.0.0)
	Address string

	// Network the server is using (default: udp)
	Network string

	// Port the server is listening on (default: 53)
	Port int

	// If the server is running in TCP mode, this listener
	// listens for incoming DNS messages
	TCPListener net.TCPListener

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

	// OOBSize is the maximum size required to store
	// out-of-band data of UDP messages. This is culculated
	// beforehand to avoid re-calculation the same size
	// over and over again
	// TODO (Techassi): Rename OOB to Ancillary
	OOBSize int

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

// New creates a new DNS server instance with the provided options
// which fallback to sane defaults
func New(c *Config) (*Server, error) {
	err := c.Validate()
	if err != nil {
		return nil, err
	}

	s := &Server{
		Address:        c.Address,
		Network:        c.Network,
		Port:           c.Port,
		UDPMessageSize: c.UDPMessageSize,
		messageList: sync.Pool{
			New: createByteBuffer(c.UDPMessageSize),
		},
		wg: sync.WaitGroup{},
	}

	err = s.init()
	return s, err
}

// createByteBuffer returns a function which creates a slice
// of bytes with the provided length
func createByteBuffer(size int) func() interface{} {
	return func() interface{} {
		return make([]byte, size)
	}
}

// createUDPListener creates a UDP listener
func createUDPListener(network, address string, port int) (*net.UDPConn, error) {
	return net.ListenUDP(network, &net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(address),
	})
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
	}

	return ErrNoSuchNetwork
}

// init initializes some server parameters
func (s *Server) init() error {
	oob4 := ipv4.NewControlMessage(ipv4.FlagDst | ipv4.FlagInterface)
	oob6 := ipv6.NewControlMessage(ipv6.FlagDst | ipv6.FlagInterface)

	if len(oob4) > len(oob6) {
		s.OOBSize = len(oob4)
	} else {
		s.OOBSize = len(oob6)
	}

	// Init new default interface implementations
	s.Resolver = resolver.NewDefaultResolver()

	s.Unpacker = pack.NewDefaultUnpacker()
	s.Packer = pack.NewDefaultPacker()

	s.Reader = reader.NewDefault(s.OOBSize)
	s.AcceptFunc = DefaultAcceptFunc

	s.Store = store.NewDefaultStore()
	s.usesCache = s.Store.UsesCache()

	err := s.Store.Set("google.de", 1, 1, net.ParseIP("142.250.186.35"))
	if err != nil {
		return err
	}

	err = s.Store.Set(".", 1, 1, net.ParseIP("198.41.0.4"))
	if err != nil {
		return err
	}

	return nil
}

// serveUDP is the main listen / respond loop of DNS messages
func (s *Server) serveUDP() error {
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
			go s.handle(m, session)
		}
	}

	return nil
}

// handle handles name matching and returns a response message
func (s *Server) handle(message dns.Message, session dns.Session) {
	// If we don't have a question just return
	if len(message.Question) == 0 {
		s.wg.Done()
		return
	}

	var record rr.RR
	var err error
	var ok bool

	// TODO (Techassi): Just pass the question instead of the individual parameters
	// First look in cache if we get a hit
	if s.usesCache {
		record, ok = s.Cache.Get(
			message.Question[0].Name,
			message.Question[0].Class,
			message.Question[0].Type,
		)
	}

	// If we don't get a cache hit or we simply use no cache
	// retrieve record data from the store
	if !ok || !s.usesCache {
		record, err = s.Store.Get(
			message.Question[0].Name,
			message.Question[0].Class,
			message.Question[0].Type,
		)
	}

	// At this point neither the cache or the store stores
	// the record data. We then try to resolve the name
	// via the resolver
	if err != nil {
		s.Resolver.Resolve(
			message.Header.RecursionDesired,
			message.Question[0].Name,
			message.Question[0].Class,
			message.Question[0].Type,
		)
	}

	message.AddAnswer(record)
	b, err := s.Packer.Pack(message)
	if err != nil {
		return
	}

	s.writeUDP(b, session)
	s.wg.Done()
}

// readUDP reads a UDP message from the UDP connection by retrieving
// a byte buffer from the message pool
func (s *Server) readUDP() ([]byte, dns.Session, error) {
	rm := s.messageList.Get().([]byte)
	mn, session, err := s.Reader.ReadUDP(s.UDPListener, rm)
	if err != nil {
		s.messageList.Put(&rm)
		return nil, session, err
	}
	return rm[:mn], session, nil
}

// writeUDPP writes a message (byte slice) back to the requesting client
func (s *Server) writeUDP(b []byte, session dns.Session) {
	_, err := s.UDPListener.WriteToUDP(b, session.Address)
	if err != nil {
		// Handle
	}
	s.messageList.Put(&b)
}

// isRunning returns if the server instance is running
func (s *Server) isRunning() bool {
	s.lock.RLock()
	running := s.running
	s.lock.RUnlock()
	return running
}

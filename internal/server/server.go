// Package server provides functions to create and manage DNS server instances
package server

import (
	"errors"
	"net"
	"sync"

	"github.com/go-void/portal/internal/cache"
	"github.com/go-void/portal/internal/config"
	"github.com/go-void/portal/internal/constants"
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
	Address net.IP

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

	// AncillarySize is the maximum size required to store
	// ancillary data of UDP messages. This is culculated
	// beforehand to avoid re-calculation the same size
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

// New creates a new DNS server instance with the provided options which fallback to sane defaults
func New(c *config.Config) (*Server, error) {
	server := &Server{
		Address:        c.Server.Address,
		Network:        c.Server.Network,
		Port:           c.Server.Port,
		UDPMessageSize: constants.UDPMinMessageSize,
		messageList: sync.Pool{
			New: createByteBuffer(constants.UDPMinMessageSize),
		},
		wg: sync.WaitGroup{},
	}

	return server, server.init()
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
	ancillary4 := ipv4.NewControlMessage(ipv4.FlagDst | ipv4.FlagInterface)
	ancillary6 := ipv6.NewControlMessage(ipv6.FlagDst | ipv6.FlagInterface)

	if len(ancillary4) > len(ancillary6) {
		s.AncillarySize = len(ancillary4)
	} else {
		s.AncillarySize = len(ancillary6)
	}

	c := cache.NewDefaultCache()

	s.Resolver = resolver.NewForwardingResolver(net.ParseIP("1.1.1.1"), c)

	s.Unpacker = pack.NewDefaultUnpacker()
	s.Packer = pack.NewDefaultPacker()

	s.Reader = reader.NewDefault(s.AncillarySize)
	s.AcceptFunc = DefaultAcceptFunc

	s.Store = store.NewDefaultStore()
	s.usesCache = s.Store.UsesCache()

	return nil
}

// createByteBuffer returns a function which creates a slice of bytes with the provided length
func createByteBuffer(size int) func() interface{} {
	return func() interface{} {
		return make([]byte, size)
	}
}

// createUDPListener creates a UDP listener
func createUDPListener(network string, address net.IP, port int) (*net.UDPConn, error) {
	return net.ListenUDP(network, &net.UDPAddr{
		Port: port,
		IP:   address,
	})
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

	var status cache.Status
	var entry cache.Entry
	var record rr.RR
	var err error

	// First look in cache if we get a hit
	if s.usesCache {
		entry, status, err = s.Cache.LookupQuestion(message.Question[0])
		record = entry.Record
	}

	// If we don't get a cache hit or we simply use no cache retrieve record data from the store
	if status != cache.Hit || !s.usesCache {
		record, err = s.Store.Get(message.Question[0])
	}

	// At this point neither the cache or the store stores
	// the record data. We then try to resolve the name
	// via the resolver
	if err != nil {
		record, err = s.Resolver.ResolveQuestion(message.Question[0])
		if err != nil {
			return
		}
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
		s.messageList.Put(rm)
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
}

// isRunning returns if the server instance is running
func (s *Server) isRunning() bool {
	s.lock.RLock()
	running := s.running
	s.lock.RUnlock()
	return running
}

// Package server provides functions to create and manage DNS server instances
package server

import (
	"errors"
	"net"
	"sync"
	"time"

	"github.com/go-void/portal/pkg/cache"
	"github.com/go-void/portal/pkg/collector"
	"github.com/go-void/portal/pkg/config"
	"github.com/go-void/portal/pkg/constants"
	"github.com/go-void/portal/pkg/dio"
	"github.com/go-void/portal/pkg/filter"
	"github.com/go-void/portal/pkg/pack"
	"github.com/go-void/portal/pkg/resolver"
	"github.com/go-void/portal/pkg/store"
	"github.com/go-void/portal/pkg/tree"
	"github.com/go-void/portal/pkg/types/dns"
	"github.com/go-void/portal/pkg/types/rr"

	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

var (
	ErrServerAlreadyRunning = errors.New("server already running")
	ErrUnexpectedConnection = errors.New("unexpected connection")
	ErrNoSuchNetwork        = errors.New("no such network")
	ErrNoQuestions          = errors.New("no questions")
)

// TODO (Techassi): Add shutdown method with context

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
	Reader dio.Reader

	// Writer implements the Writer interface to write
	// outgoing TCP and UDP messages
	Writer dio.Writer

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

	Collector collector.Collector

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

	conns sync.WaitGroup
	wg    sync.WaitGroup

	// This indicates if the server instance is running
	running bool

	// This indicates if the server / store is using a cache
	cacheEnabled bool
}

// New creates a new DNS server instance
func New() *Server {
	server := &Server{
		UDPMessageSize: constants.UDPMinMessageSize,
		messageList:    sync.Pool{},
		conns:          sync.WaitGroup{},
		wg:             sync.WaitGroup{},
	}
	server.messageList.New = createByteBuffer(constants.UDPMinMessageSize)

	return server
}

// Run starts the listen / respond loop of DNS messages
func (s *Server) Run() error {
	if s.isRunning() {
		return ErrServerAlreadyRunning
	}

	s.Collector.Run()
	s.running = true
	s.wg.Add(1)

	switch s.Network {
	case "udp", "udp4", "udp6":
		listener, err := createUDPListener(s.Network, s.Address, s.Port)
		if err != nil {
			return err
		}
		s.UDPListener = listener
		go s.serveUDP()
		return nil
	case "tcp", "tcp4", "tcp6":
		listener, err := createTCPListener(s.Network, s.Address, s.Port)
		if err != nil {
			return err
		}
		s.TCPListener = listener
		go s.serveTCP()
		return nil
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

	s.cacheEnabled = c.Server.CacheEnabled
	s.Address = c.Server.Address
	s.Network = c.Server.Network
	s.Port = c.Server.Port

	if s.Filter == nil {
		mode, err := filter.MethodFromString(c.Filter.Mode)
		if err != nil {
			mode = filter.NullMode
		}

		s.Filter = filter.New(mode)
	}

	if s.Cache == nil {
		s.Cache = cache.NewDefaultCache()
	}

	if s.Resolver == nil {
		s.Resolver = resolver.New(c.Resolver, s.Cache)
	}

	if s.Collector == nil {
		s.Collector = collector.NewDefault(c.Collector)
	}

	if s.Unpacker == nil {
		s.Unpacker = pack.NewDefaultUnpacker()
	}

	if s.Packer == nil {
		s.Packer = pack.NewDefaultPacker()
	}

	if s.Reader == nil {
		s.Reader = dio.NewDefaultReader(s.AncillarySize)
	}

	if s.Writer == nil {
		s.Writer = dio.NewDefaultWriter()
	}

	if s.AcceptFunc == nil {
		s.AcceptFunc = DefaultAcceptFunc
	}

	if s.Store == nil {
		s.Store = store.NewDefault()
	}
}

// handle handles name matching and returns a response message
func (s *Server) handle(message dns.Message, ip net.IP) (dns.Message, error) {
	start := time.Now()

	if len(message.Question) == 0 {
		s.conns.Done()
		return message, ErrNoQuestions
	}

	// TODO (Techassi): Add support for ANY queries, see RFC 8482
	var err error

	filtered, message, err := s.Filter.Match(message)
	if err != nil {
		// FIXME (Techassi): How whould we handle a filter error? Should we abort or continue (and answer the query)
		return message, err
	}

	if filtered {
		end := time.Since(start)
		entry := collector.NewFilteredEntry(message.Question[0], message.Answer[0], end, ip)
		go s.Collector.AddEntry(entry)

		return message, nil
	}

	// TODO (Techassi): Clean this up. Is there a more elegant solution?
	var status cache.Status
	var entry tree.Entry
	var record rr.RR

	if s.cacheEnabled {
		entry, status, err = s.Cache.LookupQuestion(message.Question[0])
		record = entry.Record
	}

	// TODO (Techassi): Also check if we have any custom records
	if status != cache.Hit || !s.cacheEnabled {
		record, err = s.Store.GetQuestion(message.Question[0])
	}

	if err != nil {
		record, err = s.Resolver.ResolveQuestion(message.Question[0])
		if err != nil && !errors.Is(err, resolver.ErrNoAnswer) {
			return message, err
		}
	}

	message.AddAnswer(record)
	message.IsResponse()

	end := time.Since(start)
	centry := collector.NewEntry(message.Question[0], message.Answer[0], end, ip)
	go s.Collector.AddEntry(centry)

	return message, nil
}

// isRunning returns if the server instance is running
func (s *Server) isRunning() bool {
	s.lock.RLock()
	running := s.running
	s.lock.RUnlock()
	return running
}

func (s *Server) Shutdown() {
	s.wg.Done()
}

func (s *Server) Wait() {
	s.wg.Wait()
}

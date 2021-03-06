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
	"github.com/go-void/portal/pkg/logger"
	"github.com/go-void/portal/pkg/packers"
	"github.com/go-void/portal/pkg/resolver"
	"github.com/go-void/portal/pkg/store"
	"github.com/go-void/portal/pkg/types/dns"

	"go.uber.org/zap"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

var (
	ErrServerAlreadyRunning = errors.New("server already running")
	ErrUnexpectedConnection = errors.New("unexpected connection")
	ErrNoSuchNetwork        = errors.New("no such network")
	ErrNoQuestions          = errors.New("no questions")
)

type OptionsFunc func(*Server) error

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

	// Logger is a light-weight wrapper around zap.Logger which
	// allows to the server (and all sub-components) to write
	// structured and leveled logs to one or multiple files
	Logger *logger.Logger

	// Reader implements the Reader interface to read
	// incoming TCP and UDP messages
	Reader dio.Reader

	// Writer implements the Writer interface to write
	// outgoing TCP and UDP messages
	Writer dio.Writer

	// Unpacker implements the Unpacker interface to unwrap
	// DNS messages
	Unpacker packers.Unpacker

	// Packer implements the Packer interface to pack
	// DNS messages
	Packer packers.Packer

	// Filter implements the Filter interface to enable DNS
	// filtering
	Filter filter.Engine

	// Cache implements the Cache interface to store record
	// data in memory for a specified TTL
	Cache cache.Cache

	// RecordStore implements the RecordStore interface to get and set
	// record data. A store can be anything, e.g. a file or
	// a database. The store has access to a cache instance
	// which stores record data in memory. The default store
	// does not use the cache as it is implemented as a
	// in-memory tree structure
	RecordStore store.Store

	// Resolver implements the Resolver interface to resolve
	// unknown domain names either iteratively or recursively.
	// After resolving the record data should get stored in
	// the cache
	Resolver resolver.Resolver

	// Collector implements the Collector interface to collect
	// query logs
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

	// This indicates if the server is resolving queries recursivly
	recursive bool

	config *config.Config
}

// New creates a new DNS server instance
func New(cfg *config.Config) *Server {
	server := &Server{
		UDPMessageSize: constants.UDPMinMessageSize,
		Address:        cfg.Server.Address,
		Network:        cfg.Server.Network,
		Port:           cfg.Server.Port,
		cacheEnabled:   cfg.Server.CacheEnabled,
		recursive:      cfg.Resolver.Mode == "r",
		conns:          sync.WaitGroup{},
		wg:             sync.WaitGroup{},
		messageList:    sync.Pool{},
		config:         cfg,
	}
	server.messageList.New = createByteBuffer(constants.UDPMinMessageSize)

	ancillary4 := ipv4.NewControlMessage(ipv4.FlagDst | ipv4.FlagInterface)
	ancillary6 := ipv6.NewControlMessage(ipv6.FlagDst | ipv6.FlagInterface)

	if len(ancillary4) > len(ancillary6) {
		server.AncillarySize = len(ancillary4)
	} else {
		server.AncillarySize = len(ancillary6)
	}

	return server
}

// Run starts the listen / respond loop of DNS messages
func (s *Server) Run() error {
	if s.isRunning() {
		return ErrServerAlreadyRunning
	}

	// Setup logger
	l, err := logger.New(s.config.Log)
	if err != nil {
		return err
	}
	s.Logger = l

	s.Logger.Info("start server", zap.String("context", "server"))

	// Setup defaults
	s.Defaults()

	s.Collector.Run()
	s.running = true
	s.wg.Add(1)

	switch s.Network {
	case "udp", "udp4", "udp6":
		listener, err := createUDPListener(s.Network, s.Address, s.Port)
		if err != nil {
			s.Logger.Error("failed to create UDP listener",
				zap.String("context", "server"),
				zap.Error(err),
			)
			return err
		}
		s.UDPListener = listener
		go s.serveUDP()
		return nil
	case "tcp", "tcp4", "tcp6":
		listener, err := createTCPListener(s.Network, s.Address, s.Port)
		if err != nil {
			s.Logger.Error("failed to create TCP listener",
				zap.String("context", "server"),
				zap.Error(err),
			)
			return err
		}
		s.TCPListener = listener
		go s.serveTCP()
		return nil
	}

	return ErrNoSuchNetwork
}

// Configure configures custom server components
func (s *Server) Configure(opts ...OptionsFunc) error {
	for _, opt := range opts {
		err := opt(s)
		if err != nil {
			return err
		}
	}
	return nil
}

// Defaults initializes default server parameters and checks if all neccesary
// components are registered. If not it falls back to defaults. This functions
// expects a validated config.Config
func (s *Server) Defaults() {
	if s.Filter == nil {
		s.Filter = filter.NewDefaultEngine(s.Logger)
	}

	if s.Cache == nil {
		s.Cache = cache.NewDefaultCache(s.Logger)
	}

	if s.Resolver == nil {
		s.Resolver = resolver.New(
			s.config.Resolver,
			s.Cache,
			s.Logger,
		)
	}

	if s.Collector == nil {
		s.Collector = collector.NewCollector(s.config.Collector)
	}

	if s.Unpacker == nil {
		s.Unpacker = packers.NewDefaultUnpacker()
	}

	if s.Packer == nil {
		s.Packer = packers.NewDefaultPacker()
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

	if s.RecordStore == nil {
		s.RecordStore = store.NewDefault()
	}
}

// handle handles name matching and returns a response message
func (s *Server) handle(message *dns.Message, ip net.IP) (*dns.Message, error) {
	start := time.Now()

	s.Logger.Debug("handle incoming DNS request",
		zap.String("context", "server"),
		zap.String("address", ip.String()),
		zap.Object("message", message),
	)

	if len(message.Question) == 0 {
		s.Logger.Debug("empty DNS request (no question)",
			zap.String("context", "server"),
			zap.String("address", ip.String()),
			zap.Object("message", message),
		)

		s.conns.Done()
		return message, ErrNoQuestions
	}

	// TODO (Techassi): Add support for ANY queries, see RFC 8482
	var err error

	// filtered, message, err := s.Filter.Match(ip, message)
	// if err != nil {
	// 	// FIXME (Techassi): How whould we handle a filter error? Should we abort or continue (and answer the query)
	// 	s.Logger.Error("failed to match filter",
	// 		zap.String("context", "server"),
	// 		zap.String("address", ip.String()),
	// 		zap.Object("message", message),
	// 		zap.Error(err),
	// 	)

	// 	return message, err
	// }

	// if filtered {
	// 	end := time.Since(start)
	// 	entry := collector.NewFilteredEntry(message.Question[0], message.Answer[0], end, ip)
	// 	go s.Collector.AddEntry(entry)

	// 	return message, nil
	// }

	if s.cacheEnabled {
		records, status, err := s.Cache.LookupQuestion(message.Question[0])

		if err != nil {
			return message, err
		}

		if status == cache.Hit {
			message.AddAnswers(records)
			return message, nil
		}

	}

	// Check for custom records at this point

	result, err := s.Resolver.Resolve(message)
	if err != nil {
		return message, err
	}

	// Finalize response
	message.AddRecords(result.Answer, result.Authority, result.Additional)
	message.SetIsResponse()
	message.SetRecursionAvailable(s.recursive)

	end := time.Since(start)
	centry := collector.NewEntry(message.Question[0], message.Answer, end, ip)
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
	s.Logger.Close()
	s.wg.Done()
}

func (s *Server) Wait() {
	s.wg.Wait()
}

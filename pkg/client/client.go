// Package client provides a DNS client to send and receive DNS messages
package client

import (
	"errors"
	"math"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/go-void/portal/pkg/constants"
	"github.com/go-void/portal/pkg/dio"
	"github.com/go-void/portal/pkg/logger"
	"github.com/go-void/portal/pkg/packers"
	"github.com/go-void/portal/pkg/types/dns"
	"github.com/go-void/portal/pkg/types/opcode"
	"github.com/go-void/portal/pkg/types/rcode"
	"github.com/go-void/portal/pkg/utils"

	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

var (
	ErrInvalidNetwork  = errors.New("invalid network")
	ErrNoMatchHeaderID = errors.New("header ids don't match")
)

type Client interface {
	Dial(string, net.IP) (net.Conn, error)

	// Query sends a DNS query for 'name' with 'class' and 'type'
	// to the remote DNS server with 'ip' and returns the answer
	// message and any encountered error
	Query(string, uint16, uint16, net.IP) (*dns.Message, error)

	// QueryUDP sends a DNS query using UDP
	QueryUDP(*dns.Message, net.IP) (*dns.Message, error)

	// QueryTCP sends a DNS query using TCP
	QueryTCP(*dns.Message, net.IP) (*dns.Message, error)

	// CreateMessage returns a new DNS message
	CreateMessage(dns.Header) *dns.Message

	// CreateHeader returns a new DNS message header
	// with sensible client defaults
	CreateHeader() dns.Header

	// GetID returns the current header ID and generates a new one
	GetID() uint16
}

// DefaultClient is the default DNS client to query and retrieve DNS messages
type DefaultClient struct {
	// network the client is using (default: udp)
	network string

	// unpacker implements the unpacker interface to unwrap
	// DNS messages
	unpacker packers.Unpacker

	// packer implements the packer interface to pack
	// DNS messages
	packer packers.Packer

	// reader implements the reader interface to read
	// incoming TCP and UDP messages
	reader dio.Reader

	// writer implements the writer interface to write
	// outgoing TCP and UDP messages
	writer dio.Writer

	logger *logger.Logger

	dialTimeout  time.Duration
	writeTimeout time.Duration
	readTimeout  time.Duration

	// headerID is a 16 bit uint which get's used as the DNS
	// header identifier
	headerID uint16

	lock sync.RWMutex
}

func NewDefault(l *logger.Logger) *DefaultClient {
	rand.Seed(time.Now().UnixNano())

	var size int
	ancillary4 := ipv4.NewControlMessage(ipv4.FlagDst | ipv4.FlagInterface)
	ancillary6 := ipv6.NewControlMessage(ipv6.FlagDst | ipv6.FlagInterface)

	if len(ancillary4) > len(ancillary6) {
		size = len(ancillary4)
	} else {
		size = len(ancillary6)
	}

	// TODO (Techassi): Make Client options configurable
	return &DefaultClient{
		network:      "udp",
		unpacker:     packers.NewDefaultUnpacker(),
		packer:       packers.NewDefaultPacker(),
		reader:       dio.NewDefaultReader(size),
		writer:       dio.NewDefaultWriter(),
		logger:       l,
		dialTimeout:  2 * time.Second,
		writeTimeout: 2 * time.Second,
		readTimeout:  2 * time.Second,
		headerID:     1,
	}
}

func (c *DefaultClient) Dial(network string, ip net.IP) (net.Conn, error) {
	address := utils.DNSAddress(ip)
	conn, err := net.DialTimeout(network, address, c.dialTimeout)
	if err != nil {
		return nil, err
	}

	t := time.Now()
	conn.SetWriteDeadline(t.Add(c.writeTimeout))
	conn.SetReadDeadline(t.Add(c.readTimeout))

	return conn, nil
}

// Query sends a DNS query for 'name' with 'class' and 'type' to the remote DNS server with 'ip'
// and returns the answer message and any encountered error
func (c *DefaultClient) Query(name string, class, t uint16, ip net.IP) (*dns.Message, error) {
	header := c.CreateHeader()
	query := c.CreateMessage(header)

	query.AddQuestion(dns.Question{
		Name:  name,
		Type:  t,
		Class: class,
	})

	switch c.network {
	case "udp", "udp4", "udp6":
		return c.QueryUDP(query, ip)
	}

	return nil, ErrInvalidNetwork
}

// QueryUDP sends a DNS 'query' to a remote DNS server with 'ip' using UDP
func (c *DefaultClient) QueryUDP(query *dns.Message, ip net.IP) (*dns.Message, error) {
	// Establish UDP connection
	conn, err := c.Dial(c.network, ip)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Pack DNS message into wire format
	b, err := c.packer.Pack(query)
	if err != nil {
		return nil, err
	}

	// Send query to remote DNS server
	_, err = conn.Write(b)
	if err != nil {
		return nil, err
	}

	// Read answer of the remote DNS server
	buf := make([]byte, constants.UDPMinMessageSize)
	udpConn := conn.(*net.UDPConn)
	_, _, err = c.reader.ReadUDP(udpConn, buf)
	if err != nil {
		return nil, err
	}

	// Unpack header data
	header, offset, err := c.unpacker.UnpackHeader(buf)
	if err != nil {
		return nil, err
	}

	if query.Header.ID != header.ID {
		return nil, ErrNoMatchHeaderID
	}

	// Unpack remaining message data
	return c.unpacker.Unpack(header, buf, offset)
}

// QueryTCP sends a DNS 'query' to target DNS server with 'ip' using TCP
func (c *DefaultClient) QueryTCP(query *dns.Message, ip net.IP) (*dns.Message, error) {
	conn, err := c.Dial(c.network, ip)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	b, err := c.packer.Pack(query)
	if err != nil {
		return nil, err
	}

	tcpConn := conn.(*net.TCPConn)
	err = c.writer.WriteTCPClose(tcpConn, b)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// CreateMessage creates a new DNS message with a header
func (c *DefaultClient) CreateMessage(header dns.Header) *dns.Message {
	return &dns.Message{
		Header: header,
	}
}

// CreateHeader returns a new DNS message header
func (c *DefaultClient) CreateHeader() dns.Header {
	return dns.Header{
		ID:                 c.GetID(),
		IsQuery:            true,
		OpCode:             opcode.Query,
		Authoritative:      false,
		Truncated:          false,
		RecursionDesired:   true,
		RecursionAvailable: false,
		Zero:               false,
		RCode:              rcode.NoError,
	}
}

// GetID returns the current header ID and generates a new one
func (c *DefaultClient) GetID() uint16 {
	c.lock.Lock()
	id := c.headerID
	c.headerID = uint16(rand.Intn(math.MaxUint16-1) + 1)
	c.lock.Unlock()
	return id
}

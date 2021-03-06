// Package client provides a DNS client to send and receive DNS messages
package client

import (
	"errors"
	"math"
	"math/rand"
	"net"
	"net/netip"
	"sync"
	"time"

	"github.com/go-void/portal/pkg/constants"
	"github.com/go-void/portal/pkg/dio"
	"github.com/go-void/portal/pkg/logger"
	"github.com/go-void/portal/pkg/packers"
	"github.com/go-void/portal/pkg/types/dns"

	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

var (
	ErrInvalidNetwork  = errors.New("invalid network")
	ErrNoMatchHeaderID = errors.New("header ids don't match")
)

// Client is the default DNS client to query and retrieve DNS messages
type Client struct {
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

func New(l *logger.Logger) *Client {
	rand.Seed(time.Now().UnixNano())

	var size int
	ancillary4 := ipv4.NewControlMessage(ipv4.FlagDst | ipv4.FlagInterface)
	ancillary6 := ipv6.NewControlMessage(ipv6.FlagDst | ipv6.FlagInterface)

	if len(ancillary4) > len(ancillary6) {
		size = len(ancillary4)
	} else {
		size = len(ancillary6)
	}

	return &Client{
		network:      "udp",
		unpacker:     packers.NewDefaultUnpacker(),
		packer:       packers.NewDefaultPacker(),
		reader:       dio.NewDefaultReader(size),
		writer:       dio.NewDefaultWriter(),
		dialTimeout:  2 * time.Second,
		writeTimeout: 2 * time.Second,
		readTimeout:  2 * time.Second,
		headerID:     1,
		logger:       l,
	}
}

// Configure configures available client options
func (c *Client) Configure(opts ...OptionFunc) error {
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) Dial(network string, addr netip.AddrPort) (net.Conn, error) {
	conn, err := net.DialTimeout(network, addr.String(), c.dialTimeout)
	if err != nil {
		return nil, err
	}

	t := time.Now()
	conn.SetWriteDeadline(t.Add(c.writeTimeout))
	conn.SetReadDeadline(t.Add(c.readTimeout))

	return conn, nil
}

// Query sends a DNS query for 'name' with 'class' and 'type' to the remote DNS server with 'ip' and returns the answer
// message and any encountered error
func (c *Client) Query(name string, class, t uint16, addr netip.Addr) (*dns.Message, error) {
	header := dns.NewHeader(c.GetID())
	query := dns.NewMessageWith(header)

	query.AddQuestion(dns.Question{
		Name:  name,
		Type:  t,
		Class: class,
	})

	switch c.network {
	case "udp", "udp4", "udp6":
		return c.QueryUDP(query, netip.AddrPortFrom(addr, 53))
	}

	return nil, ErrInvalidNetwork
}

// QueryUDP sends a DNS 'query' to a remote DNS server with 'ip' using UDP
func (c *Client) QueryUDP(query *dns.Message, addrPort netip.AddrPort) (*dns.Message, error) {
	// Establish UDP connection
	conn, err := c.Dial(c.network, addrPort)
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
func (c *Client) QueryTCP(query *dns.Message, addrPort netip.AddrPort) (*dns.Message, error) {
	conn, err := c.Dial(c.network, addrPort)
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

// GetID returns the current header ID and generates a new one
func (c *Client) GetID() uint16 {
	c.lock.Lock()
	id := c.headerID
	c.headerID = uint16(rand.Intn(math.MaxUint16-1) + 1)
	c.lock.Unlock()
	return id
}

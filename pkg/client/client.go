// Package client provides a DNS client to send and receive DNS messages
package client

import (
	"encoding/binary"
	"errors"
	"math"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/go-void/portal/pkg/pack"
	"github.com/go-void/portal/pkg/reader"
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
	// Query sends a DNS query for 'name' with 'class' and 'type'
	// to the remote DNS server with 'ip' and returns the answer
	// message and any encountered error
	Query(string, uint16, uint16, net.IP) (dns.Message, error)

	// QueryUDP sends a DNS query using UDP
	QueryUDP(dns.Message, net.IP) (dns.Message, error)

	// QueryTCP sends a DNS query using TCP
	QueryTCP(dns.Message, net.IP) (dns.Message, error)

	// CreateMessage returns a new DNS message
	CreateMessage(dns.Header) dns.Message

	// CreateHeader returns a new DNS message header
	// with sensible client defaults
	CreateHeader() dns.Header

	// GetID returns the current header ID and generates a new one
	GetID() uint16
}

// DefaultClient is the default DNS client to query and retrieve DNS messages
type DefaultClient struct {
	// Network the client is using (default: udp)
	Network string

	// Unpacker implements the Unpacker interface to unwrap
	// DNS messages
	Unpacker pack.Unpacker

	// Packer implements the Packer interface to pack
	// DNS messages
	Packer pack.Packer

	// Reader implements the Reader interface to read
	// incoming TCP and UDP messages
	Reader reader.Reader

	DialTimeout  time.Duration
	WriteTimeout time.Duration
	ReadTimeout  time.Duration

	// headerID is a 16 bit uint which get's used as the DNS
	// header identifier
	headerID uint16

	lock sync.RWMutex
}

func NewDefault() *DefaultClient {
	rand.Seed(time.Now().UnixNano())

	var size = 0
	ancillary4 := ipv4.NewControlMessage(ipv4.FlagDst | ipv4.FlagInterface)
	ancillary6 := ipv6.NewControlMessage(ipv6.FlagDst | ipv6.FlagInterface)

	if len(ancillary4) > len(ancillary6) {
		size = len(ancillary4)
	} else {
		size = len(ancillary6)
	}

	// TODO (Techassi): Make Client options configurable
	return &DefaultClient{
		Network:      "udp",
		Unpacker:     pack.NewDefaultUnpacker(),
		Packer:       pack.NewDefaultPacker(),
		Reader:       reader.NewDefault(size),
		DialTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		ReadTimeout:  2 * time.Second,
		headerID:     1,
	}
}

// Query sends a DNS query for 'name' with 'class' and 'type' to the remote DNS server with 'ip'
// and returns the answer message and any encountered error
func (c *DefaultClient) Query(name string, class, t uint16, ip net.IP) (dns.Message, error) {
	header := c.CreateHeader()
	query := c.CreateMessage(header)

	query.AddQuestion(dns.Question{
		Name:  name,
		Type:  t,
		Class: class,
	})

	switch c.Network {
	case "udp", "udp4", "udp6":
		return c.QueryUDP(query, ip)
	}

	return dns.Message{}, ErrInvalidNetwork
}

// QueryUDP sends a DNS query using UDP
func (c *DefaultClient) QueryUDP(query dns.Message, ip net.IP) (dns.Message, error) {
	address := utils.DNSAddress(ip)
	conn, err := net.DialTimeout("udp", address, c.DialTimeout)
	if err != nil {
		return dns.Message{}, err
	}
	defer conn.Close()

	t := time.Now()
	conn.SetWriteDeadline(t.Add(c.WriteTimeout))
	conn.SetReadDeadline(t.Add(c.ReadTimeout))

	b, err := c.Packer.Pack(query)
	if err != nil {
		return dns.Message{}, err
	}

	_, err = conn.Write(b)
	if err != nil {
		return dns.Message{}, err
	}

	buf := make([]byte, 512)
	c.Reader.ReadUDP(&conn, buf)

	header, offset, err := c.Unpacker.UnpackHeader(buf)
	if err != nil {
		return dns.Message{}, err
	}

	if query.Header.ID != header.ID {
		return dns.Message{}, ErrNoMatchHeaderID
	}

	return c.Unpacker.Unpack(header, buf, offset)
}

// QueryTCP sends a DNS query using TCP
func (c *DefaultClient) QueryTCP(query dns.Message, ip net.IP) (dns.Message, error) {
	address := utils.DNSAddress(ip)
	conn, err := net.DialTimeout("tcp", address, c.DialTimeout)
	if err != nil {
		return dns.Message{}, err
	}
	defer conn.Close()

	t := time.Now()
	conn.SetWriteDeadline(t.Add(c.WriteTimeout))
	conn.SetReadDeadline(t.Add(c.ReadTimeout))

	b, err := c.Packer.Pack(query)
	if err != nil {
		return dns.Message{}, err
	}

	q := make([]byte, len(b)+2)
	binary.BigEndian.PutUint16(q, uint16(len(b)))
	copy(q[2:], b)

	_, err = conn.Write(q)
	if err != nil {
		return dns.Message{}, err
	}

	return dns.Message{}, nil
}

// CreateMessage creates a new DNS message with a header
func (c *DefaultClient) CreateMessage(header dns.Header) dns.Message {
	return dns.Message{
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

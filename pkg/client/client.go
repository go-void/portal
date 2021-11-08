// Package client provides a DNS client to send and receive DNS messages
package client

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/go-void/portal/pkg/pack"
	"github.com/go-void/portal/pkg/types/dns"
	"github.com/go-void/portal/pkg/types/opcode"
	"github.com/go-void/portal/pkg/types/rcode"
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

	// QueryTCO sends a DNS query using TCP
	QueryTCP(dns.Message, net.IP) (dns.Message, error)

	// CreateMessage returns a new DNS message
	CreateMessage(dns.Header) dns.Message

	// CreateHeader returns a new DNS message header
	// with sensible client defaults
	CreateHeader() dns.Header

	// GetID returns the current header ID and increments
	// it by one
	GetID() uint16
}

type DefaultClient struct {
	// Network the client is using (default: udp)
	Network string

	// Unpacker implements the Unpacker interface to unwrap
	// DNS messages
	Unpacker pack.Unpacker

	// Packer implements the Packer interface to pack
	// DNS messages
	Packer pack.Packer

	DialTimeout  time.Duration
	WriteTimeout time.Duration
	ReadTimeout  time.Duration

	// headerID is a 16 bit uint which get's used as the DNS
	// header identifier
	headerID uint16

	lock sync.RWMutex
}

func NewDefaultClient() *DefaultClient {
	rand.Seed(time.Now().UnixNano())

	// TODO (Techassi): Make Client options configurable
	return &DefaultClient{
		Network:      "udp",
		Unpacker:     pack.NewDefaultUnpacker(),
		Packer:       pack.NewDefaultPacker(),
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
	var address string

	if len(ip) == 16 {
		address = fmt.Sprintf("%s:%d", ip.String(), 53)
	} else {
		address = fmt.Sprintf("[%s]:%d", ip.String(), 53)
	}

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

	var buf = make([]byte, 512)
	offset, err := conn.Read(buf)
	if err != nil {
		return dns.Message{}, err
	}
	buf = buf[:offset]

	header, offset, err := c.Unpacker.UnpackHeader(buf)
	if err != nil {
		return dns.Message{}, err
	}

	if query.Header.ID != header.ID {
		return dns.Message{}, ErrNoMatchHeaderID
	}

	return c.Unpacker.Unpack(header, buf, offset)
}

// QueryTCO sends a DNS query using TCP
func (c *DefaultClient) QueryTCP(query dns.Message, ip net.IP) (dns.Message, error) {
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

// GetID returns the current header ID and increments it by one
func (c *DefaultClient) GetID() uint16 {
	c.lock.Lock()
	id := c.headerID
	c.headerID = uint16(rand.Intn(math.MaxUint16-1) + 1)
	c.lock.Unlock()
	return id
}

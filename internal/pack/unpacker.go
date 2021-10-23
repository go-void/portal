package pack

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strings"

	"github.com/go-void/portal/internal/types/dns"
	"github.com/go-void/portal/internal/types/rr"
	"github.com/go-void/portal/internal/utils"
)

var (
	ErrNoBody          = errors.New("no body data")
	ErrUnpackpingQName = errors.New("error while unpacking QNAME")
)

type Unpacker interface {
	// Unpack unwraps a single complete DNS message from the
	// received byte slice
	Unpack(dns.MessageHeader, []byte, int) (dns.Message, error)

	// UnpackHeader unwraps header data from the received
	// byte slice
	UnpackHeader([]byte) (dns.MessageHeader, int, error)

	// UnpackQuestion unwraps a question from the received
	// byte slice
	UnpackQuestion([]byte, int) (dns.Question, int, error)

	// UnpackName unwraps a domain name in a DNS question
	// or in a RR header
	UnpackName([]byte, int) (string, int, error)

	// UnpackRRList unwraps a list of resource records from the
	// received byte slice
	UnpackRRList(uint16, []byte, int) ([]rr.RR, int, error)

	// UnpackRR unwraps a single resource record from the
	// received byte slice
	UnpackRR([]byte, int) (rr.RR, int, error)

	// UnpackRRHeader unwraps header data of a resource
	// record from the received byte slice
	UnpackRRHeader([]byte, int) (rr.RRHeader, int, error)
}

// DefaultWrapper describes the default wrapper to unwrap / wrap
// DNS messages. It is based on a regular byte Reader to read
// bytes into a pre-defined struct
type DefaultUnpacker struct {
	reader  *bytes.Reader
	builder strings.Builder
}

// NewDefaultUnpacker creates a new default wrapper instance
func NewDefaultUnpacker() Unpacker {
	return &DefaultUnpacker{
		reader: bytes.NewReader([]byte{}),
	}
}

// Unpack unwraps a single complete DNS message from the received byte slice
func (p *DefaultUnpacker) Unpack(header dns.MessageHeader, data []byte, offset int) (dns.Message, error) {
	m := dns.Message{}
	m.Header = header

	// Immediatly return if the message only consists of header data
	// without any body data
	if offset == len(data) {
		return m, ErrNoBody
	}

	// We cannot trust the values of QDCOUNT, ANCOUNT, NSCOUNT and
	// ARCOUNT, as these values can be manipulated by potential
	// attackers. The first step is to assume the values are correct
	// and if we detect a wrong offset we can be pretty sure the
	// count is wrong

	// Loop over the questions. Usually there is only one question,
	// but the spec accounts for the possibility to ask multiple
	// questions at once
	for i := 0; i < int(header.QDCount); i++ {
		// Save initial offset to compare later
		initialOffset := offset

		question, o, err := p.UnpackQuestion(data, offset)
		if err != nil {
			return m, err
		}
		offset = o

		// If the initial offset and the offset after unwrapping
		// the question match we know that QDCOUNT is wrong
		if initialOffset == o {
			header.QDCount = uint16(i)
			break
		}

		m.Question = append(m.Question, question)
	}

	// Unpack slice of answer RRS in the answer section
	answers, offset, err := p.UnpackRRList(header.ANCount, data, offset)
	if err != nil {
		return m, err
	}
	m.Answer = answers

	// Unpack slice of nameserver RRs in the authority section
	nameservers, offset, err := p.UnpackRRList(header.NSCount, data, offset)
	if err != nil {
		return m, err
	}
	m.Authority = nameservers

	// Unpack slice of additional RRs in the additional section
	additional, _, err := p.UnpackRRList(header.ARCount, data, offset)
	if err != nil {
		return m, err
	}
	m.Additional = additional

	return m, nil
}

// UnpackHeader unwraps header data from the received byte slice
func (p *DefaultUnpacker) UnpackHeader(data []byte) (dns.MessageHeader, int, error) {
	p.reader.Reset(data)

	rh := new(dns.RawHeader)
	err := binary.Read(p.reader, binary.BigEndian, rh)
	if err != nil {
		return dns.MessageHeader{}, 0, err
	}

	return rh.ToHeader(), binary.Size(rh), nil
}

// UnpackQuestion unwraps a question from the received byte slice
func (p *DefaultUnpacker) UnpackQuestion(data []byte, offset int) (dns.Question, int, error) {
	qname, offset, err := p.UnpackName(data, offset)
	if err != nil {
		return dns.Question{}, offset, err
	}

	t, offset := utils.UnpackUint16(data, offset)
	c, offset := utils.UnpackUint16(data, offset)

	q := dns.Question{
		Name:  qname,
		Type:  t,
		Class: c,
	}

	return q, offset, nil
}

// UnpackName unwraps a domain name in a DNS question or in a RR header
func (p *DefaultUnpacker) UnpackName(data []byte, offset int) (string, int, error) {
	// TODO (Techassi): Optimize this

	// If we immediation encounter a null byte, the name is root (.)
	if data[offset] == 0x00 {
		return ".", offset + 1, nil
	}

	// Initialize the end of the label in bytes
	end := offset + int(data[offset]) + 1
	p.builder.Reset()
	offset++

	// Iterate over the bytes until we reach the null byte, which
	// marks the root (.)
	for {
		if data[offset] == 0x00 {
			_, err := p.builder.WriteString(".")
			if err != nil {
				return p.builder.String(), offset, ErrUnpackpingQName
			}

			offset++
			break
		}

		if offset == end {
			_, err := p.builder.WriteString(".")
			if err != nil {
				return p.builder.String(), offset, ErrUnpackpingQName
			}

			end += int(data[offset]) + 1
			offset++
		}

		err := p.builder.WriteByte(data[offset])
		if err != nil {
			return p.builder.String(), offset, ErrUnpackpingQName
		}

		offset++
	}

	return p.builder.String(), offset, nil
}

// UnpackRRList unwraps a list of resource records from the received byte slice
func (p *DefaultUnpacker) UnpackRRList(count uint16, data []byte, offset int) ([]rr.RR, int, error) {
	if count == 0 {
		return nil, offset, nil
	}

	var list []rr.RR

	for i := 0; i < int(count); i++ {
		initialOffset := offset
		rr, o, err := p.UnpackRR(data, offset)
		if err != nil {
			continue
		}
		offset = o

		// If the initial offset and the offset after unwrapping
		// the RR match we know that count is wrong
		if initialOffset == o {
			break
		}
		list = append(list, rr)
	}

	return list, offset, nil
}

// UnpackRR unwraps a single resource record from the received byte slice
func (p *DefaultUnpacker) UnpackRR(data []byte, offset int) (rr.RR, int, error) {
	header, offset, err := p.UnpackRRHeader(data, offset)
	if err != nil {
		return nil, offset, err
	}

	record, err := rr.New(header.Type)
	if err != nil {
		return nil, offset, err
	}

	// TODO (Techassi): Check RDLENGTH
	record.SetHeader(header)

	offset, err = record.Unpack(data, offset)
	if err != nil {
		return nil, offset, err
	}

	return record, offset, nil
}

// UnpackRRHeader unwraps header data of a resource record from the received byte slice
func (p *DefaultUnpacker) UnpackRRHeader(data []byte, offset int) (rr.RRHeader, int, error) {
	header := rr.RRHeader{}

	// Unpack NAME
	name, offset, err := p.UnpackName(data, offset)
	if err != nil {
		return header, offset, err
	}
	header.Name = name

	// Unpack TYPE
	rrType, offset := utils.UnpackUint16(data, offset)
	header.Type = rrType

	// Unpack CLASS
	rrClass, offset := utils.UnpackUint16(data, offset)
	header.Class = rrClass

	// Unpack TTL
	rrTTL, offset := utils.UnpackUint32(data, offset)
	header.TTL = rrTTL

	// Unpack RDLENGTH
	rdlength, offset := utils.UnpackUint16(data, offset)
	header.RDLength = rdlength

	return header, offset, nil
}

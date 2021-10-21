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
	ErrUnwrappingQName = errors.New("error while unwrapping QNAME")
)

type Unpacker interface {
	// Unpack unwraps a single complete DNS message from the
	// received byte slice
	Unpack(dns.MessageHeader, []byte, int) (dns.Message, error)

	// UnwrapHeader unwraps header data from the received
	// byte slice
	UnpackHeader([]byte) (dns.MessageHeader, int, error)

	// UnpackQuestion unwraps a question from the received
	// byte slice
	UnpackQuestion([]byte, int) (dns.Question, int, error)

	// UnpackName unwraps a domain name in a DNS question
	// or in a RR header
	UnpackName([]byte, int) (string, int, error)

	// UnwrapRRList unwraps a list of resource records from the
	// received byte slice
	UnpackRRList(uint16, []byte, int) ([]rr.RR, int, error)

	// UnwrapRR unwraps a single resource record from the
	// received byte slice
	UnpackRR([]byte, int) (rr.RR, int, error)

	// UnwrapRRHeader unwraps header data of a resource
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

// Unwrap unwraps a single complete DNS message from the received byte slice
func (w *DefaultUnpacker) Unpack(header dns.MessageHeader, data []byte, offset int) (dns.Message, error) {
	m := dns.Message{}

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

		question, o, err := w.UnpackQuestion(data, offset)
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

	// Unwrap slice of answer RRS in the answer section
	answers, offset, err := w.UnpackRRList(header.ANCount, data, offset)
	if err != nil {
		return m, err
	}
	m.Answer = answers

	// Unwrap slice of nameserver RRs in the authority section
	nameservers, offset, err := w.UnpackRRList(header.NSCount, data, offset)
	if err != nil {
		return m, err
	}
	m.Authority = nameservers

	// Unwrap slice of additional RRs in the additional section
	additional, _, err := w.UnpackRRList(header.ARCount, data, offset)
	if err != nil {
		return m, err
	}
	m.Additional = additional

	return m, nil
}

// UnwrapHeader unwraps header data from the received byte slice
func (w *DefaultUnpacker) UnpackHeader(data []byte) (dns.MessageHeader, int, error) {
	w.reader.Reset(data)

	rh := new(dns.RawHeader)
	err := binary.Read(w.reader, binary.BigEndian, rh)
	if err != nil {
		return dns.MessageHeader{}, 0, err
	}

	return rh.ToHeader(), binary.Size(rh), nil
}

// UnpackQuestion unwraps a question from the received byte slice
func (w *DefaultUnpacker) UnpackQuestion(data []byte, offset int) (dns.Question, int, error) {
	qname, offset, err := w.UnpackName(data, offset)
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
func (w *DefaultUnpacker) UnpackName(data []byte, offset int) (string, int, error) {
	// TODO (Techassi): Optimize this

	// If we immediation encounter a null byte, the name is root (.)
	if data[offset] == 0x00 {
		return ".", offset + 1, nil
	}

	// Initialize the end of the label in bytes
	end := offset + int(data[offset]) + 1
	w.builder.Reset()
	offset++

	// Iterate over the bytes until we reach the null byte, which
	// marks the root (.)
	for {
		if data[offset] == 0x00 {
			_, err := w.builder.WriteString(".")
			if err != nil {
				return w.builder.String(), offset, ErrUnwrappingQName
			}

			offset++
			break
		}

		if offset == end {
			_, err := w.builder.WriteString(".")
			if err != nil {
				return w.builder.String(), offset, ErrUnwrappingQName
			}

			end += int(data[offset]) + 1
			offset++
		}

		err := w.builder.WriteByte(data[offset])
		if err != nil {
			return w.builder.String(), offset, ErrUnwrappingQName
		}

		offset++
	}

	return w.builder.String(), offset, nil
}

// UnwrapRRList unwraps a list of resource records from the received byte slice
func (w *DefaultUnpacker) UnpackRRList(count uint16, data []byte, offset int) ([]rr.RR, int, error) {
	if count == 0 {
		return nil, offset, nil
	}

	var list []rr.RR

	for i := 0; i < int(count); i++ {
		initialOffset := offset
		rr, o, err := w.UnpackRR(data, offset)
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

// UnwrapRR unwraps a single resource record from the received byte slice
func (w *DefaultUnpacker) UnpackRR(data []byte, offset int) (rr.RR, int, error) {
	header, offset, err := w.UnpackRRHeader(data, offset)
	if err != nil {
		return nil, offset, err
	}

	record, err := rr.New(header.Type)
	if err != nil {
		return nil, offset, err
	}

	// TODO (Techassi): Check RDLENGTH
	record.SetHeader(header)

	offset, err = record.Unwrap(data, offset)
	if err != nil {
		return nil, offset, err
	}

	return record, offset, nil
}

// UnwrapRRHeader unwraps header data of a resource record from the received byte slice
func (w *DefaultUnpacker) UnpackRRHeader(data []byte, offset int) (rr.RRHeader, int, error) {
	header := rr.RRHeader{}

	// Unwrap NAME
	name, offset, err := w.UnpackName(data, offset)
	if err != nil {
		return header, offset, err
	}
	header.Name = name

	// Unwrap TYPE
	rrType, offset := utils.UnpackUint16(data, offset)
	header.Type = rrType

	// Unwrap CLASS
	rrClass, offset := utils.UnpackUint16(data, offset)
	header.Class = rrClass

	// Unwrap TTL
	rrTTL, offset := utils.UnpackUint32(data, offset)
	header.TTL = rrTTL

	// Unwrap RDLENGTH
	rdlength, offset := utils.UnpackUint16(data, offset)
	header.RDLength = rdlength

	return header, offset, nil
}

package rr

// This sectiom implements all standard RRs mentioned in 3.3.
// Standard RRs.
// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.3

// This sectiom implements all standard RRs mentioned in 3.4.
// Internet specific RRs.
// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.4

import (
	"errors"
)

var (
	ErrNoSuchType            = errors.New("no such RR type")
	ErrFailedToConvertRRDate = errors.New("failed to convert RR data")
	ErrInvalidRRData         = errors.New("invalid RR data")
)

// RR describes a resource record. Information of a (domain) name
// is composed of a set of resource records.
// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.2
type RR interface {
	// Get header of RR
	Header() *RRHeader

	// Set header of RR
	SetHeader(RRHeader)

	// Set resource record data
	SetData(...interface{}) error

	// String returns the representation of any RR as text
	String() string

	// Unwrap unwraps the RDATA
	Unwrap([]byte, int) (int, error)
}

// RRHeader describes header data of a resource record.
// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.2.1
type RRHeader struct {
	// Name specifies an owner name
	Name string

	// Type specifies a RR type code
	Type uint16

	// Class specifies a RR class code
	Class uint16

	// TTL specifies the time interval in seconds that the RR
	// may be cached before it should be considered outdated
	TTL uint32

	// RDLength specifies the length of RDATA in octets
	RDLength uint16
}

// NOTE (Techassi): Maybe generate this
var typeMap = map[uint16]func() RR{
	TypeCNAME: func() RR { return new(CNAME) },
	TypeHINFO: func() RR { return new(HINFO) },
	TypeMB:    func() RR { return new(MB) },
	TypeMD:    func() RR { return new(MD) },
	TypeMF:    func() RR { return new(MF) },
	TypeMG:    func() RR { return new(MG) },
	TypeMINFO: func() RR { return new(MINFO) },
	TypeMR:    func() RR { return new(MR) },
	TypeMX:    func() RR { return new(MX) },
	TypeNULL:  func() RR { return new(NULL) },
	TypeNS:    func() RR { return new(NS) },
	TypePTR:   func() RR { return new(PTR) },
	TypeSOA:   func() RR { return new(SOA) },
	TypeTXT:   func() RR { return new(TXT) },
	TypeA:     func() RR { return new(A) },
	TypeWKS:   func() RR { return new(WKS) },
}

// New returns a new RR based on the provided type
func New(t uint16) (RR, error) {
	rr, ok := typeMap[t]
	if !ok {
		return nil, ErrNoSuchType
	}

	return rr(), nil
}

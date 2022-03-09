package rr

// This sectiom implements all standard RRs mentioned in 3.3.
// Standard RRs.
// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.3

// This sectiom implements all standard RRs mentioned in 3.4.
// Internet specific RRs.
// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.4

import (
	"errors"
	"math"
	"strings"
	"time"

	"github.com/go-void/portal/pkg/compression"
)

var (
	ErrNoSuchType            = errors.New("no such RR type")
	ErrFailedToConvertRRData = errors.New("failed to convert RR data")
	ErrInvalidRRData         = errors.New("invalid RR data")
)

// RR describes a resource record. Information of a (domain) name
// is composed of a set of resource records.
// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.2
type RR interface {
	// Get header of RR
	Header() *Header

	// Set header of RR
	SetHeader(Header)

	// Set resource record data
	SetData(...interface{}) error

	// String returns the representation of any RR as text
	String() string

	// Len returns the records RDLENGTH
	Len() uint16

	// IsSame compares both RRs and returns if they are the same
	IsSame(RR) bool

	// Unpack unpacks the RDATA
	Unpack([]byte, int) (int, error)

	// Pack packs the RDATA
	Pack([]byte, int, compression.Map) (int, error)
}

// Header describes header data of a resource record.
// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.2.1
type Header struct {
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

	// Additional Expires field. This helps to more easily handle cached records
	Expires int64
}

// New returns a new RR based on the provided type
func New(t uint16) (RR, error) {
	create, ok := typeMap[t]
	if !ok {
		return nil, ErrNoSuchType
	}

	return create(), nil
}

// NewFromName returns a new RR based on the provided name
func NewFromName(name string) (RR, uint16, error) {
	t := stringToTypeMap[strings.ToUpper(name)]
	rr, err := New(t)
	if err != nil {
		return nil, 0, err
	}

	return rr, t, nil
}

// UpdateTTL updates the TTL of a record based on the expiry timestamp
func UpdateTTL(record RR, expire time.Time) {
	h := record.Header()
	ttl := uint32(math.Max(0, (time.Until(expire)).Seconds()))
	h.TTL = ttl
}

package dns

import (
	"github.com/go-void/portal/pkg/types/opcode"
	"github.com/go-void/portal/pkg/types/rcode"

	m "github.com/go-void/portal/pkg/types/bitmasks"
)

// Header describes the header data of a message. This header format
// enables easy access to all header fields. The RawHeader in
// comparison stores raw data directly from the "wire".
// See https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.1
type Header struct {
	ID                 uint16      // ID
	IsQuery            bool        // QR
	OpCode             opcode.Code // OPCODE
	Authoritative      bool        // AA
	Truncated          bool        // TC
	RecursionDesired   bool        // RD
	RecursionAvailable bool        // RA
	Zero               bool        // Z
	RCode              rcode.Code  // RCODE
	QDCount            uint16      // Question count
	ANCount            uint16      // Answer count
	NSCount            uint16      // Authority count
	ARCount            uint16      // Additional record count
}

// RawHeader describes the raw header data of a message
// directly from the "wire". The data gets unpacked by
// splitting the message into six 16 bit (2 octet)
// chunks. The second chunk "flags" carries data like
// QR, OPCODE, etc. which gets split up further by bit
// masks
type RawHeader struct {
	ID      uint16 // ID
	Flags   uint16 // Various flags, see above
	QDCount uint16 // QDCOUNT
	ANCount uint16 // ANCOUNT
	NSCount uint16 // NSCOUNT
	ARCount uint16 // ARCOUNT
}

func NewHeader(id uint16) Header {
	return Header{
		ID:                 id,
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

// ToHeader converts a raw header to a header by applying
// bitmasks to split DNS header flags
func (h *RawHeader) ToHeader() Header {
	return Header{
		ID:                 h.ID,
		IsQuery:            h.Flags&m.QR == 0,
		OpCode:             opcode.Code((h.Flags >> 11) & 0xF),
		Authoritative:      h.Flags&m.AA != 0,
		Truncated:          h.Flags&m.TC != 0,
		RecursionDesired:   h.Flags&m.RD != 0,
		RecursionAvailable: h.Flags&m.RA != 0,
		Zero:               h.Flags&m.Z != 0,
		RCode:              rcode.Code(h.Flags & 0xF),
		QDCount:            h.QDCount,
		ANCount:            h.ANCount,
		NSCount:            h.NSCount,
		ARCount:            h.ARCount,
	}
}

// ToRaw converts a header to a raw header by applying
// bitmasks to shift data to the correct positions
func (h *Header) ToRaw() RawHeader {
	var rh RawHeader

	rh.ID = h.ID
	rh.Flags = uint16(h.OpCode)<<11 | uint16(h.RCode&0xF)

	if !h.IsQuery {
		rh.Flags |= m.QR
	}

	if h.Authoritative {
		rh.Flags |= m.AA
	}

	if h.Truncated {
		rh.Flags |= m.TC
	}

	if h.RecursionDesired {
		rh.Flags |= m.RD
	}

	if h.RecursionAvailable {
		rh.Flags |= m.RA
	}

	if h.Zero {
		rh.Flags |= m.Z
	}

	rh.QDCount = h.QDCount
	rh.ANCount = h.ANCount
	rh.NSCount = h.NSCount
	rh.ARCount = h.ARCount

	return rh
}

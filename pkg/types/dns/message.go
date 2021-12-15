package dns

import (
	"github.com/go-void/portal/pkg/constants"
	"github.com/go-void/portal/pkg/labels"
	"github.com/go-void/portal/pkg/types/opcode"
	"github.com/go-void/portal/pkg/types/rcode"
	"github.com/go-void/portal/pkg/types/rr"

	m "github.com/go-void/portal/pkg/types/bitmasks"
)

// Message describes a complete DNS message describes in RFC 1035
// Section 4.
// See https://datatracker.ietf.org/doc/html/rfc1035#section-4
type Message struct {
	Header     Header
	Question   []Question
	Answer     []rr.RR
	Authority  []rr.RR
	Additional []rr.RR

	// Compression keeps track of compression pointers
	// and domain names
	Compression CompressionMap
}

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

func NewMessage() Message {
	return Message{
		Compression: NewCompressionMap(),
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

func (m *Message) SetIsResponse() {
	m.Header.IsQuery = false
}

func (m *Message) SetRecursionAvailable(ra bool) {
	m.Header.RecursionAvailable = m.Header.RecursionDesired && ra
}

// AddQuestion adds a question to the question section
// of a DNS message
func (m *Message) AddQuestion(question Question) {
	m.Question = append(m.Question, question)
	m.Header.QDCount++
}

// AddAnswer adds a resource record to the answer section
// of a DNS message
func (m *Message) AddAnswer(record rr.RR) {
	if record == nil {
		return
	}

	m.Answer = append(m.Answer, record)
	m.Header.ANCount++
}

func (m *Message) AddAnswers(records []rr.RR) {
	if records == nil {
		return
	}

	m.Answer = append(m.Answer, records...)
	m.Header.ANCount += uint16(len(records))
}

// AddAuthority adds a resource record to the
// authoritative name server section
func (m *Message) AddAuthority(record rr.RR) {
	if record == nil {
		return
	}

	m.Answer = append(m.Authority, record)
	m.Header.NSCount++
}

// AddAdditional adds a resource record to the
// additional section
func (m *Message) AddAdditional(record rr.RR) {
	if record == nil {
		return
	}

	m.Additional = append(m.Additional, record)
	m.Header.ARCount++
}

func (m *Message) Len() int {
	// Fixed DNS header length
	len := constants.DNSHeaderLen

	// DNS question length
	len += labels.Len(m.Question[0].Name)
	len += constants.DNSQuestionFixedLen

	for _, a := range m.Answer {
		len += int(a.Len())
	}

	for _, a := range m.Authority {
		len += int(a.Len())
	}

	for _, a := range m.Additional {
		len += int(a.Len())
	}

	return len
}

// IsEDNS returns if the message has an EDNS OPT record
func (m *Message) IsEDNS() bool {
	// We iterate from the back because the OPT RR is usually at the
	// end of the additional records
	for i := len(m.Additional) - 1; i >= 0; i-- {
		if _, ok := m.Additional[i].(*rr.OPT); ok {
			return true
		}
	}
	return false
}

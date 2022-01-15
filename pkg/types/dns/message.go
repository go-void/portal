package dns

import (
	"github.com/go-void/portal/pkg/constants"
	"github.com/go-void/portal/pkg/labels"
	"github.com/go-void/portal/pkg/types/rr"

	"go.uber.org/zap/zapcore"
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

func NewMessage() *Message {
	return &Message{
		Compression: NewCompressionMap(),
	}
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

// MarshalLogObject marshals the DNS message as a zap log object
func (m Message) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddUint16("id", m.Header.ID)
	return enc.AddObject("question", m.Question[0])
}

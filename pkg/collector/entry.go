package collector

import (
	"net"
	"time"

	"github.com/go-void/portal/pkg/types/dns"
	"github.com/go-void/portal/pkg/types/rr"
)

type Entry struct {
	ID            string
	Question      dns.Question
	Answer        []rr.RR
	QueryTime     time.Duration
	ClientIP      net.IP
	AppliedFilter string
	Filtered      bool
	Cached        bool
}

func NewEntry(question dns.Question, answer []rr.RR, queryTime time.Duration, ip net.IP) Entry {
	return Entry{
		ID:            "",
		QueryTime:     queryTime,
		Question:      question,
		Answer:        answer,
		ClientIP:      ip,
		AppliedFilter: "",
		Filtered:      false,
		Cached:        false,
	}
}

func NewCachedEntry(question dns.Question, answer []rr.RR, queryTime time.Duration, ip net.IP) Entry {
	entry := NewEntry(question, answer, queryTime, ip)
	entry.Cached = true
	return entry
}

func NewFilteredEntry(question dns.Question, answer []rr.RR, queryTime time.Duration, ip net.IP) Entry {
	entry := NewEntry(question, answer, queryTime, ip)
	entry.Filtered = true
	return entry
}

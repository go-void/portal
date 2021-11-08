// Package filter provides different filters to filter out DNS requests (to block them)
package filter

import (
	"errors"
	"net"
	"strings"

	"github.com/go-void/portal/pkg/labels"
	"github.com/go-void/portal/pkg/types/dns"
	"github.com/go-void/portal/pkg/types/rcode"
	"github.com/go-void/portal/pkg/types/rr"
)

var (
	ErrInvalidFilterMethod = errors.New("filter: invalid filter method")
	ErrInvalidDomainRule   = errors.New("filter: invalid domain rule")
	ErrInvalidIPAddress    = errors.New("filter: invalid ip address")
	ErrNoSuchRule          = errors.New("filter: no such rule")

	defaultFilterIP = net.IPv4(0, 0, 0, 0)
)

type RuleType int

const (
	DomainRule RuleType = iota
	RPZRule
)

type Filter interface {
	Match(dns.Message) (bool, dns.Message, error)

	ParseRule(RuleType, string) (string, net.IP, error)

	AddRule(RuleType, string) error

	RemoveRule(string) error

	Mode() FilterMode
}

// TODO (Techassi): Make configurable
func New(mode FilterMode) Filter {
	return &DefaultFilter{
		Rules:      make(map[string]net.IP),
		FilterMode: mode,
		TTL:        300,
	}
}

type DefaultFilter struct {
	// FilterMode defines how the filter should answer a filtered request
	FilterMode FilterMode

	// TTL defines the TTL (in seconds) returned by filtered answers
	TTL int

	// Address defines the IP address of this DNS server
	Address net.IP

	// Rules stores a map of rules
	Rules map[string]net.IP
}

func (f *DefaultFilter) Match(message dns.Message) (bool, dns.Message, error) {
	var question = message.Question[0]

	if ip, ok := f.Rules[question.Name]; ok {
		switch f.FilterMode {
		case NxDomainMode:
			message.Header.RCode = rcode.NameError
			return true, message, nil
		case LocalIPMode:
			answer, err := rr.New(question.Type)
			if err != nil {
				return true, message, err
			}

			err = answer.SetData(f.Address)
			if err != nil {
				return true, message, err
			}

			answer.SetHeader(rr.Header{
				Name:     question.Name,
				Type:     question.Type,
				Class:    question.Class,
				TTL:      uint32(f.TTL),
				RDLength: answer.Len(),
			})

			message.AddAnswer(answer)
			return true, message, nil
		case NoDataMode:
			return true, message, nil
		case NullMode:
			answer, err := rr.New(question.Type)
			if err != nil {
				return true, message, err
			}

			err = answer.SetData(ip)
			if err != nil {
				return true, message, err
			}

			answer.SetHeader(rr.Header{
				Name:     question.Name,
				Type:     question.Type,
				Class:    question.Class,
				TTL:      uint32(f.TTL),
				RDLength: answer.Len(),
			})

			message.AddAnswer(answer)
			return true, message, nil
		}
	}

	return false, message, nil
}

// ParseRule parses a filter rule with the following format: '<ip-address> <domain>'.
// Example: '0.0.0.0 example.com'
func (f *DefaultFilter) ParseRule(t RuleType, input string) (string, net.IP, error) {
	// TODO (Techassi): Check domain validity (labels.IsValidDomain)
	// FIXME (Techassi): Actually check rule type
	parts := strings.Split(input, " ")
	switch len(parts) {
	case 1:
		return parts[0], defaultFilterIP, nil
	case 2:
		ip := net.ParseIP(parts[0])
		if ip == nil {
			return "", nil, ErrInvalidIPAddress
		}

		return parts[1], ip, nil
	default:
		return "", nil, ErrInvalidDomainRule
	}
}

func (f *DefaultFilter) AddRule(t RuleType, rule string) error {
	domain, ip, err := f.ParseRule(t, rule)
	if err != nil {
		return err
	}

	domain = labels.Rootify(domain)
	f.Rules[domain] = ip

	return nil
}

func (f *DefaultFilter) RemoveRule(domain string) error {
	if _, ok := f.Rules[domain]; ok {
		delete(f.Rules, domain)
		return nil
	}
	return ErrNoSuchRule
}

func (f *DefaultFilter) Mode() FilterMode {
	return f.FilterMode
}

// Package filter provides different filters to filter out DNS requests (to block them)
package filter

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/go-void/portal/internal/labels"
	"github.com/go-void/portal/internal/types/dns"
	"github.com/go-void/portal/internal/types/rcode"
	"github.com/go-void/portal/internal/types/rr"
)

var (
	ErrInvalidFilterMethod = errors.New("filter: invalid filter method")
	ErrInvalidDomainMatch  = errors.New("filter: tried to match invalid domain")
	ErrInvalidDomainRule   = errors.New("filter: invalid domain rule")
	ErrInvalidIPAddress    = errors.New("filter: invalid ip address")

	defaultFilterIP = net.IPv4(0, 0, 0, 0)
)

type RuleType int

const (
	DomainRule RuleType = iota
	RPZRule
)

type Filter interface {
	ParseRule(string) (string, net.IP, error)

	AddRulesFromFile(RuleType, string) error

	AddRulesFromURL(RuleType, string) error

	AddRule(RuleType, string) error

	Match(dns.Message) (bool, dns.Message, error)

	Mode() FilterMode
}

// TODO (Techassi): Make configurable
func New() Filter {
	return &DefaultFilter{
		Rules:      make(map[string]net.IP),
		FilterMode: NullMode,
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

	// WatchURLS indicates if URL watching is enabled
	WatchURLS bool

	// WatchFiles indicates if file watching is enabled
	WatchFiles bool

	// WatchInterval defines the watch interval in which URLs should
	// be checked for changes
	WatchInterval int

	// Rules stores a map of rules
	Rules map[string]net.IP

	// Files stores a slice of file paths
	Files []string

	// Urls stores a slice of URLs
	Urls []string
}

// ParseRule parses a filter rule with the following format: '<ip-address> <domain>'.
// Example: '0.0.0.0 example.com'
func (f *DefaultFilter) ParseRule(input string) (string, net.IP, error) {
	// TODO (Techassi): Check domain validity (labels.IsValidDomain)
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

func (f *DefaultFilter) AddRulesFromFile(t RuleType, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	s := bufio.NewScanner(file)

	for s.Scan() {
		domain, ip, err := f.ParseRule(s.Text())
		if err != nil {
			return err
		}

		domain = labels.Rootify(domain)
		f.Files = append(f.Files, path)
		f.Rules[domain] = ip
	}

	return nil
}

func (f *DefaultFilter) AddRulesFromURL(t RuleType, url string) error {
	// TODO (Techassi): Don't use default HTTP client, use custom one so that the user can adjust options
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("filter: request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("filter: failed to read body: %v", err)
	}

	r := bytes.NewReader(body)
	s := bufio.NewScanner(r)

	for s.Scan() {
		domain, ip, err := f.ParseRule(s.Text())
		if err != nil {
			return err
		}

		domain = labels.Rootify(domain)
		f.Urls = append(f.Urls, url)
		f.Rules[domain] = ip
	}

	return nil
}

func (f *DefaultFilter) AddRule(t RuleType, rule string) error {
	// FIXME (Techassi): Actually check rule type
	domain, ip, err := f.ParseRule(rule)
	if err != nil {
		return err
	}

	domain = labels.Rootify(domain)
	f.Rules[domain] = ip

	return nil
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

func (f *DefaultFilter) Mode() FilterMode {
	return f.FilterMode
}

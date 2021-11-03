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
)

var (
	ErrInvalidDomainMatch = errors.New("filter: tried to match invalid domain")
	ErrInvalidDomainRule  = errors.New("filter: invalid domain rule")
)

type DefaultFilter struct {
	filterMethod FilterMethod
	ttl          int

	rules map[string]net.IP
	files []string
	urls  []string
}

// ParseRule parses a filter rule with the following format: '<ip-address> <domain>'.
// Example: '0.0.0.0 example.com'
func (f *DefaultFilter) ParseRule(input string) (string, net.IP, error) {
	// FIXME (Techassi): Also support a rule without an IP address and then default back to 0.0.0.0

	parts := strings.Split(input, " ")
	if len(parts) != 2 {
		return "", nil, ErrInvalidDomainRule
	}

	ip := net.ParseIP(parts[0])
	if ip == nil {
		return "", nil, ErrInvalidIPAddress
	}

	// TODO (Techassi): Check domain validity (labels.IsValidDomain)
	return parts[1], ip, nil
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

		f.files = append(f.files, path)
		f.rules[domain] = ip
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

		f.urls = append(f.urls, url)
		f.rules[domain] = ip
	}

	return nil
}

func (f *DefaultFilter) AddRule(t RuleType, rule string) error {
	// FIXME (Techassi): Actually check rule type
	domain, ip, err := f.ParseRule(rule)
	if err != nil {
		return err
	}

	f.rules[domain] = ip
	return nil
}

func (f *DefaultFilter) Match(domain string) (FilterResult, error) {
	// FIXME (Techassi): Check if the domain ends with a '.'.
	// Example: Rule 0.0.0.0 example.com should also filter example.com.

	if !labels.IsValidDomain(domain) {
		return FilterResult{}, ErrInvalidDomainMatch
	}

	if ip, ok := f.rules[domain]; ok {
		return FilterResult{
			Filtered: true,
			Rule:     domain,
			Target:   ip,
		}, nil
	}

	return FilterResult{}, nil
}

func (f *DefaultFilter) Method() FilterMethod {
	return f.filterMethod
}

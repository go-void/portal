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
)

var (
	ErrInvalidDomainRule = errors.New("filter: invalid domain rule")
)

type DomainFilter struct {
	rules map[string]net.IP
	files []string
	urls  []string
}

// ParseRule parses a filter rule with the following format: '<ip-address> <domain>'.
// Example: '0.0.0.0 example.com'
func (f *DomainFilter) ParseRule(input string) (string, net.IP, error) {
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

// LoadFromString loads a set of filter rules from a string input. Rules have to be separated by newlines
func (f *DomainFilter) LoadFromString(input string) error {
	r := strings.NewReader(input)
	s := bufio.NewScanner(r)

	for s.Scan() {
		domain, ip, err := f.ParseRule(s.Text())
		if err != nil {
			return err
		}

		f.rules[domain] = ip
	}

	return nil
}

// LoadFromFile loads a set of filter rules from a file. Rules have to be separated by newlines
func (f *DomainFilter) LoadFromFile(path string) error {
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

// LoadFromURL loads a set of filter rules from URL. Rules have to be separated by newlines
func (f *DomainFilter) LoadFromURL(url string) error {
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

func (f *DomainFilter) Refresh() error {
	panic("not implemented") // TODO: Implement
}

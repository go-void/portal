// Package filter provides different filters to filter out DNS requests (to block them)
package filter

import (
	"errors"
	"net"
)

var (
	ErrInvalidFilterMethod = errors.New("filter: invalid filter method")
	ErrInvalidIPAddress    = errors.New("filter: invalid ip address")
)

type RuleType int

const (
	DomainRule RuleType = iota
	RPZRule
)

type FilterResult struct {
	Filtered bool
	Rule     string
	Target   net.IP
}

type Filter interface {
	ParseRule(string) (string, net.IP, error)

	AddRulesFromFile(RuleType, string) error

	AddRulesFromURL(RuleType, string) error

	AddRule(RuleType, string) error

	Match(string) (FilterResult, error)

	Method() FilterMethod
}

// TODO (Techassi): Make configurable
func New() Filter {
	return &DefaultFilter{
		rules:        make(map[string]net.IP),
		filterMethod: NullMethod,
		ttl:          300,
	}
}

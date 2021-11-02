// Package filter provides different filters to filter out DNS requests (to block them)
package filter

import (
	"errors"
	"net"
)

var (
	ErrInvalidIPAddress = errors.New("filter: invalid ip address")
)

type Filter interface {
	// ParseRule parses a filter rule
	ParseRule(string) (string, net.IP, error)

	LoadFromString(string) error

	LoadFromFile(string) error

	LoadFromURL(string) error

	Refresh() error
}

func NewDomainFilter() Filter {
	return &DomainFilter{}
}

func NewRPZFilter() Filter {
	return &RPZFilter{}
}

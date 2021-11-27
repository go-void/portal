// Package filter provides different functions to filter out DNS requests (to block them)
package filter

import (
	"errors"
	"net"

	"github.com/go-void/portal/pkg/types/dns"
)

var (
	ErrNoSuchFilter = errors.New("filter: no such filter")
)

type Engine interface {
	AddFilter(Filter) error

	GetFilter(*net.IP) (*Filter, error)

	Match(net.IP, dns.Message) (bool, dns.Message, error)
}

type DefaultEngine struct {
	ipToFilter map[*net.IP]int
	filters    []Filter
}

func NewEngine() Engine {
	return &DefaultEngine{}
}

func (e *DefaultEngine) AddFilter(f Filter) error {
	// Code ...
	return nil
}

func (e *DefaultEngine) GetFilter(*net.IP) (*Filter, error) {
	return nil, ErrNoSuchFilter
}

func (e *DefaultEngine) Match(ip net.IP, message dns.Message) (bool, dns.Message, error) {
	filter, err := e.GetFilter(&ip)
	if err != nil {
		return false, message, nil
	}

	return filter.Match(message)
}

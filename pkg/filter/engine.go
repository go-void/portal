// Package filter provides different functions to filter out DNS requests (to block them)
package filter

import "net"

type Engine interface {
	AddFilter(Filter) error

	GetFilter(*net.IP) (*Filter, error)
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
	return nil, nil
}

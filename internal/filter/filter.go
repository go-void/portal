// Package filter provides different filters to filter out DNS requests (to block them)
package filter

type Filter interface {
}

func NewHostFilter() Filter {
	return &HostFilter{}
}

func NewRPZFilter() Filter {
	return &RPZFilter{}
}

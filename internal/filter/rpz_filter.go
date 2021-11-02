package filter

import "net"

type RPZFilter struct {
}

func (f *RPZFilter) ParseRule(input string) (string, net.IP, error) {
	return "", nil, nil
}

func (f *RPZFilter) LoadFromString(_ string) error {
	panic("not implemented") // TODO: Implement
}

func (f *RPZFilter) LoadFromFile(_ string) error {
	panic("not implemented") // TODO: Implement
}

func (f *RPZFilter) LoadFromURL(_ string) error {
	panic("not implemented") // TODO: Implement
}

func (f *RPZFilter) Refresh() error {
	panic("not implemented") // TODO: Implement
}

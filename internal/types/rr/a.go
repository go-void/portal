package rr

import (
	"net"
)

// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.4.1
type A struct {
	H       RRHeader
	Address net.IP
}

func (rr *A) Header() *RRHeader {
	return &rr.H
}

func (rr *A) SetHeader(header RRHeader) {
	rr.H = header
}

func (rr *A) SetData(data ...interface{}) error {
	if len(data) != 1 {
		return ErrInvalidRRData
	}

	addr, ok := data[0].(net.IP)
	if !ok {
		return ErrFailedToConvertRRDate
	}
	rr.Address = addr
	return nil
}

func (rr *A) String() string {
	return rr.Address.String()
}

func (rr *A) Unwrap(data []byte, offset int) (int, error) {
	return offset, nil
}

package rr

import "net"

// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.4.2
type WKS struct {
	H        RRHeader
	Address  net.IP
	Protocol [1]byte
	BitMap   []byte
}

func (rr *WKS) Header() *RRHeader {
	return &rr.H
}

func (rr *WKS) SetHeader(header RRHeader) {
	rr.H = header
}

func (rr *WKS) SetData(data ...interface{}) error {
	if len(data) != 3 {
		return ErrInvalidRRData
	}

	// TODO (Techassi): Add fields

	return nil
}

func (rr *WKS) String() string {
	return ""
}

func (rr *WKS) Unwrap(data []byte, offset int) (int, error) {
	return offset, nil
}

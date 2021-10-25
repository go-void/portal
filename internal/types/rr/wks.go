package rr

import (
	"net"

	"github.com/go-void/portal/internal/wire"
)

// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.4.2
type WKS struct {
	H        Header
	Address  net.IP
	Protocol [1]byte
	BitMap   []byte
}

func (rr *WKS) Header() *Header {
	return &rr.H
}

func (rr *WKS) SetHeader(header Header) {
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

func (rr *WKS) Unpack(data []byte, offset int) (int, error) {
	return offset, nil
}

func (rr *WKS) Pack(buf []byte, offset int) (int, error) {
	offset, err := wire.PackIPAddress(rr.Address, buf, offset)
	if err != nil {
		return offset, err
	}

	offset, err = wire.PackUint8(rr.Protocol[0], buf, offset)
	if err != nil {
		return offset, err
	}

	return wire.PackBitMap(rr.BitMap, buf, offset)
}

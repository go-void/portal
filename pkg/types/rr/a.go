package rr

import (
	"fmt"
	"net/netip"

	"github.com/go-void/portal/pkg/compression"
	"github.com/go-void/portal/pkg/pack"
)

// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.4.1
type A struct {
	H       Header
	Address netip.Addr
}

func (rr *A) Header() *Header {
	return &rr.H
}

func (rr *A) SetHeader(header Header) {
	rr.H = header
}

func (rr *A) SetData(data ...interface{}) error {
	if len(data) != 1 {
		return ErrInvalidRRData
	}

	addr, ok := data[0].(netip.Addr)
	if !ok {
		return ErrFailedToConvertRRData
	}
	rr.Address = addr
	return nil
}

func (rr *A) String() string {
	return fmt.Sprintf("A <%v %s>", rr.H, rr.Address)
}

func (rr *A) Len() uint16 {
	return 4
}

func (rr *A) IsSame(o RR) bool {
	other, ok := o.(*A)
	if !ok {
		return false
	}

	return rr.Address == other.Address
}

func (rr *A) Unpack(data []byte, offset int) (int, error) {
	address, offset, err := pack.UnpackIPv4Address(data, offset)
	if err != nil {
		return offset, err
	}
	rr.Address = address

	return offset, nil
}

func (rr *A) Pack(buf []byte, offset int, _ compression.Map) (int, error) {
	return pack.PackIPAddress(rr.Address, buf, offset)
}

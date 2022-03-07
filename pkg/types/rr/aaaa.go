package rr

import (
	"fmt"
	"net"

	"github.com/go-void/portal/pkg/pack"
)

// See https://datatracker.ietf.org/doc/html/rfc3596
type AAAA struct {
	H       Header
	Address net.IP
}

func (rr *AAAA) Header() *Header {
	return &rr.H
}

func (rr *AAAA) SetHeader(header Header) {
	rr.H = header
}

func (rr *AAAA) SetData(data ...interface{}) error {
	if len(data) != 1 {
		return ErrInvalidRRData
	}

	addr, ok := data[0].(net.IP)
	if !ok {
		return ErrFailedToConvertRRData
	}
	rr.Address = addr
	return nil
}

func (rr *AAAA) String() string {
	return fmt.Sprintf("%v %s", rr.H, rr.Address)
}

func (rr *AAAA) Len() uint16 {
	return 16
}

func (rr *AAAA) IsSame(o RR) bool {
	other, ok := o.(*AAAA)
	if !ok {
		return false
	}

	return rr.Address.Equal(other.Address)
}

func (rr *AAAA) Unpack(data []byte, offset int) (int, error) {
	address, offset, err := pack.UnpackIPv6Address(data, offset)
	if err != nil {
		return offset, err
	}
	rr.Address = address

	return offset, nil
}

func (rr *AAAA) Pack(buf []byte, offset int) (int, error) {
	return pack.PackIPAddress(rr.Address, buf, offset)
}

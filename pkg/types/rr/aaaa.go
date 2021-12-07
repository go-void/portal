package rr

import (
	"encoding/binary"
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

func (rr *AAAA) Unpack(data []byte, offset int) (int, error) {
	address, offset := pack.UnpackIPv6Address(data, offset)
	rr.Address = address
	return offset, nil
}

func (rr *AAAA) Pack(buf []byte, offset int) (int, error) {
	ip := binary.BigEndian.Uint32(rr.Address[12:16])
	binary.BigEndian.PutUint32(buf[offset:], ip)
	return offset + 4, nil
}

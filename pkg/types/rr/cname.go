package rr

import (
	"github.com/go-void/portal/pkg/compression"
	"github.com/go-void/portal/pkg/pack"
)

// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.3.1
type CNAME struct {
	H      Header
	Target string
}

func (rr *CNAME) Header() *Header {
	return &rr.H
}

func (rr *CNAME) SetHeader(header Header) {
	rr.H = header
}

func (rr *CNAME) SetData(data ...interface{}) error {
	if len(data) != 1 {
		return ErrInvalidRRData
	}

	target, ok := data[0].(string)
	if !ok {
		return ErrFailedToConvertRRData
	}
	rr.Target = target
	return nil
}

func (rr *CNAME) String() string {
	return ""
}

func (rr *CNAME) Len() uint16 {
	return uint16(len(rr.Target)) + 1
}

func (rr *CNAME) IsSame(o RR) bool {
	other, ok := o.(*CNAME)
	if !ok {
		return false
	}

	return rr.Target == other.Target
}

func (rr *CNAME) Unpack(data []byte, offset int) (int, error) {
	target, offset, err := pack.UnpackDomainName(data, offset)
	if err != nil {
		return offset, err
	}
	rr.Target = target

	return offset, nil
}

func (rr *CNAME) Pack(buf []byte, offset int, comp compression.Map) (int, error) {
	return pack.PackDomainName(rr.Target, buf, offset, comp)
}

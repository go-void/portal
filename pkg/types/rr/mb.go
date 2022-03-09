package rr

import (
	"github.com/go-void/portal/pkg/compression"
	"github.com/go-void/portal/pkg/pack"
)

// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.3.3 (EXPERIMENTAL)
type MB struct {
	H       Header
	MADName string
}

func (rr *MB) Header() *Header {
	return &rr.H
}

func (rr *MB) SetHeader(header Header) {
	rr.H = header
}

func (rr *MB) SetData(data ...interface{}) error {
	if len(data) != 1 {
		return ErrInvalidRRData
	}

	name, ok := data[0].(string)
	if !ok {
		return ErrFailedToConvertRRData
	}
	rr.MADName = name

	return nil
}

func (rr *MB) String() string {
	return ""
}

func (rr *MB) Len() uint16 {
	return uint16(len(rr.MADName)) + 1
}

func (rr *MB) IsSame(o RR) bool {
	other, ok := o.(*MB)
	if !ok {
		return false
	}

	return rr.MADName == other.MADName
}

func (rr *MB) Unpack(data []byte, offset int) (int, error) {
	name, offset, err := pack.UnpackDomainName(data, offset)
	if err != nil {
		return offset, err
	}
	rr.MADName = name

	return offset, nil
}

func (rr *MB) Pack(buf []byte, offset int, comp compression.Map) (int, error) {
	return pack.PackDomainName(rr.MADName, buf, offset, comp)
}

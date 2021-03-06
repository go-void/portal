package rr

import (
	"github.com/go-void/portal/pkg/compression"
	"github.com/go-void/portal/pkg/pack"
)

// https://datatracker.ietf.org/doc/html/rfc1035#section-3.3.4 (Obsolete)
type MD struct {
	H       Header
	MADName string
}

func (rr *MD) Header() *Header {
	return &rr.H
}

func (rr *MD) SetHeader(header Header) {
	rr.H = header
}

func (rr *MD) SetData(data ...interface{}) error {
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

func (rr *MD) String() string {
	return ""
}

func (rr *MD) Len() uint16 {
	return uint16(len(rr.MADName)) + 1
}

func (rr *MD) IsSame(o RR) bool {
	other, ok := o.(*MD)
	if !ok {
		return false
	}

	return rr.MADName == other.MADName
}

func (rr *MD) Unpack(data []byte, offset int) (int, error) {
	name, offset, err := pack.UnpackDomainName(data, offset)
	if err != nil {
		return offset, err
	}
	rr.MADName = name

	return offset, nil
}

func (rr *MD) Pack(buf []byte, offset int, comp compression.Map) (int, error) {
	return pack.PackDomainName(rr.MADName, buf, offset, comp)
}

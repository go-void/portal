package rr

import "github.com/go-void/portal/pkg/pack"

// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.3.12
type PTR struct {
	H        Header
	PTRDName string
}

func (rr *PTR) Header() *Header {
	return &rr.H
}

func (rr *PTR) SetHeader(header Header) {
	rr.H = header
}

func (rr *PTR) SetData(data ...interface{}) error {
	if len(data) != 1 {
		return ErrInvalidRRData
	}

	name, ok := data[0].(string)
	if !ok {
		return ErrFailedToConvertRRData
	}
	rr.PTRDName = name
	return nil
}

func (rr *PTR) String() string {
	return ""
}

func (rr *PTR) Len() uint16 {
	return uint16(len(rr.PTRDName)) + 1
}

func (rr *PTR) IsSame(o RR) bool {
	other, ok := o.(*PTR)
	if !ok {
		return false
	}

	return rr.PTRDName == other.PTRDName
}

func (rr *PTR) Unpack(data []byte, offset int) (int, error) {
	pointer, offset, err := pack.UnpackDomainName(data, offset)
	if err != nil {
		return offset, err
	}
	rr.PTRDName = pointer

	return offset, nil
}

func (rr *PTR) Pack(buf []byte, offset int) (int, error) {
	return pack.PackDomainName(rr.PTRDName, buf, offset)
}

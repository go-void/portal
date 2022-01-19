package rr

import "github.com/go-void/portal/pkg/pack"

// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.3.8 (EXPERIMENTAL)
type MR struct {
	H       Header
	NewName string
}

func (rr *MR) Header() *Header {
	return &rr.H
}

func (rr *MR) SetHeader(header Header) {
	rr.H = header
}

func (rr *MR) SetData(data ...interface{}) error {
	if len(data) != 1 {
		return ErrInvalidRRData
	}

	name, ok := data[0].(string)
	if !ok {
		return ErrFailedToConvertRRData
	}
	rr.NewName = name

	return nil
}

func (rr *MR) String() string {
	return ""
}

func (rr *MR) Len() uint16 {
	return uint16(len(rr.NewName)) + 1
}

func (rr *MR) IsSame(o RR) bool {
	other, ok := o.(*MR)
	if !ok {
		return false
	}

	return rr.NewName == other.NewName
}

func (rr *MR) Unpack(data []byte, offset int) (int, error) {
	name, offset := pack.UnpackDomainName(data, offset)
	rr.NewName = name
	return offset, nil
}

func (rr *MR) Pack(buf []byte, offset int) (int, error) {
	return pack.PackDomainName(rr.NewName, buf, offset)
}

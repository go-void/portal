package rr

import "github.com/go-void/portal/pkg/pack"

// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.3.2
type HINFO struct {
	H   Header
	CPU string
	OS  string
}

func (rr *HINFO) Header() *Header {
	return &rr.H
}

func (rr *HINFO) SetHeader(header Header) {
	rr.H = header
}

func (rr *HINFO) SetData(data ...interface{}) error {
	if len(data) != 2 {
		return ErrInvalidRRData
	}

	cpu, ok := data[0].(string)
	if !ok {
		return ErrFailedToConvertRRData
	}
	rr.CPU = cpu

	os, ok := data[1].(string)
	if !ok {
		return ErrFailedToConvertRRData
	}
	rr.OS = os

	return nil
}

func (rr *HINFO) String() string {
	return ""
}

func (rr *HINFO) Len() uint16 {
	return uint16(len(rr.CPU)+len(rr.OS)) + 2
}

func (rr *HINFO) IsSame(o RR) bool {
	other, ok := o.(*HINFO)
	if !ok {
		return false
	}

	return rr.CPU == other.CPU && rr.OS == other.OS
}

func (rr *HINFO) Unpack(data []byte, offset int) (int, error) {
	cpu, offset, err := pack.UnpackCharacterString(data, offset)
	if err != nil {
		return offset, err
	}
	rr.CPU = cpu

	os, offset, err := pack.UnpackCharacterString(data, offset)
	if err != nil {
		return offset, err
	}
	rr.OS = os

	return offset, nil
}

func (rr *HINFO) Pack(buf []byte, offset int) (int, error) {
	offset, err := pack.PackCharacterString(rr.CPU, buf, offset)
	if err != nil {
		return offset, err
	}

	return pack.PackCharacterString(rr.OS, buf, offset)
}

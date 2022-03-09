package rr

import (
	"github.com/go-void/portal/pkg/compression"
	"github.com/go-void/portal/pkg/pack"
)

// TODO (Techassi): Figure out if a TXT record can hold multiple strings

// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.3.14
type TXT struct {
	H    Header
	Data string
}

func (rr *TXT) Header() *Header {
	return &rr.H
}

func (rr *TXT) SetHeader(header Header) {
	rr.H = header
}

func (rr *TXT) SetData(data ...interface{}) error {
	if len(data) != 1 {
		return ErrInvalidRRData
	}

	d, ok := data[0].(string)
	if !ok {
		return ErrFailedToConvertRRData
	}
	rr.Data = d
	return nil
}

func (rr *TXT) String() string {
	return ""
}

func (rr *TXT) Len() uint16 {
	return uint16(len(rr.Data)) + 1
}

func (rr *TXT) IsSame(o RR) bool {
	other, ok := o.(*TXT)
	if !ok {
		return false
	}

	return rr.Data == other.Data
}

func (rr *TXT) Unpack(data []byte, offset int) (int, error) {
	// TODO (Techassi): Figure out how to unpack multiple character strings
	str, offset, err := pack.UnpackCharacterString(data, offset)
	if err != nil {
		return offset, err
	}
	rr.Data = str

	return offset, nil
}

func (rr *TXT) Pack(buf []byte, offset int, _ compression.Map) (int, error) {
	return pack.PackCharacterString(rr.Data, buf, offset)
}

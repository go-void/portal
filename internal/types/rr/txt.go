package rr

import "github.com/go-void/portal/internal/wire"

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

func (rr *TXT) Unpack(data []byte, offset int) (int, error) {
	// TODO (Techassi): Figure out how to unpack multiple character strings
	return offset, nil
}

func (rr *TXT) Pack(buf []byte, offset int) (int, error) {
	return wire.PackCharacterString(rr.Data, buf, offset)
}

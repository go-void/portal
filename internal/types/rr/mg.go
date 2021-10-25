package rr

import "github.com/go-void/portal/internal/wire"

// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.3.6 (EXPERIMENTAL)
type MG struct {
	H       Header
	MGMName string
}

func (rr *MG) Header() *Header {
	return &rr.H
}

func (rr *MG) SetHeader(header Header) {
	rr.H = header
}

func (rr *MG) SetData(data ...interface{}) error {
	if len(data) != 1 {
		return ErrInvalidRRData
	}

	name, ok := data[0].(string)
	if !ok {
		return ErrFailedToConvertRRDate
	}
	rr.MGMName = name

	return nil
}

func (rr *MG) String() string {
	return ""
}

func (rr *MG) Unpack(data []byte, offset int) (int, error) {
	return offset, nil
}

func (rr *MG) Pack(buf []byte, offset int) (int, error) {
	return wire.PackDomainName(rr.MGMName, buf, offset)
}

package rr

import "github.com/go-void/portal/internal/wire"

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
		return ErrFailedToConvertRRDate
	}
	rr.MADName = name

	return nil
}

func (rr *MD) String() string {
	return ""
}

func (rr *MD) Unpack(data []byte, offset int) (int, error) {
	return offset, nil
}

func (rr *MD) Pack(buf []byte, offset int) (int, error) {
	return wire.PackDomainName(rr.MADName, buf, offset)
}

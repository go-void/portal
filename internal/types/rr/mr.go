package rr

import "github.com/go-void/portal/internal/wire"

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

func (rr *MR) Unpack(data []byte, offset int) (int, error) {
	name, offset := wire.UnpackDomainName(data, offset)
	rr.NewName = name
	return offset, nil
}

func (rr *MR) Pack(buf []byte, offset int) (int, error) {
	return wire.PackDomainName(rr.NewName, buf, offset)
}

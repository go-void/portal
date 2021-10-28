package rr

import (
	"fmt"

	"github.com/go-void/portal/internal/wire"
)

// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.3.11
type NS struct {
	H       Header
	NSDName string
}

func (rr *NS) Header() *Header {
	return &rr.H
}

func (rr *NS) SetHeader(header Header) {
	rr.H = header
}

func (rr *NS) SetData(data ...interface{}) error {
	if len(data) != 1 {
		return ErrInvalidRRData
	}

	name, ok := data[0].(string)
	if !ok {
		return ErrFailedToConvertRRDate
	}
	rr.NSDName = name
	return nil
}

func (rr *NS) String() string {
	return fmt.Sprintf("%v %s", rr.H, rr.NSDName)
}

func (rr *NS) Unpack(data []byte, offset int) (int, error) {
	name, offset := wire.UnpackDomainName(data, offset)
	rr.NSDName = name
	return offset, nil
}

func (rr *NS) Pack(buf []byte, offset int) (int, error) {
	return wire.PackDomainName(rr.NSDName, buf, offset)
}

package rr

import (
	"fmt"

	"github.com/go-void/portal/pkg/compression"
	"github.com/go-void/portal/pkg/pack"
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
		return ErrFailedToConvertRRData
	}
	rr.NSDName = name
	return nil
}

func (rr *NS) String() string {
	return fmt.Sprintf("%v %s", rr.H, rr.NSDName)
}

func (rr *NS) Len() uint16 {
	return uint16(len(rr.NSDName)) + 1
}

func (rr *NS) IsSame(o RR) bool {
	other, ok := o.(*NS)
	if !ok {
		return false
	}

	return rr.NSDName == other.NSDName
}

func (rr *NS) Unpack(data []byte, offset int) (int, error) {
	name, offset, err := pack.UnpackDomainName(data, offset)
	if err != nil {
		return offset, err
	}
	rr.NSDName = name

	return offset, nil
}

func (rr *NS) Pack(buf []byte, offset int, comp compression.Map) (int, error) {
	return pack.PackDomainName(rr.NSDName, buf, offset, comp)
}

package rr

import (
	"fmt"

	"github.com/go-void/portal/pkg/compression"
	"github.com/go-void/portal/pkg/labels"
	"github.com/go-void/portal/pkg/pack"
)

// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.3.9
type MX struct {
	H          Header
	Preference uint16
	Exchange   string
}

func (rr *MX) Header() *Header {
	return &rr.H
}

func (rr *MX) SetHeader(header Header) {
	rr.H = header
}

func (rr *MX) SetData(data ...interface{}) error {
	if len(data) != 2 {
		return ErrInvalidRRData
	}

	pref, ok := data[0].(uint16)
	if !ok {
		return ErrFailedToConvertRRData
	}
	rr.Preference = pref

	exchange, ok := data[1].(string)
	if !ok {
		return ErrFailedToConvertRRData
	}
	rr.Exchange = exchange

	return nil
}

func (rr *MX) String() string {
	return fmt.Sprintf("MX <%v Preference: %d, Exchange: %s>", rr.H, rr.Preference, rr.Exchange)
}

func (rr *MX) Len() uint16 {
	return uint16(labels.Len(rr.Exchange)) + 2
}

func (rr *MX) IsSame(o RR) bool {
	other, ok := o.(*MX)
	if !ok {
		return false
	}

	return rr.Preference == other.Preference && rr.Exchange == other.Exchange
}

func (rr *MX) Unpack(data []byte, offset int) (int, error) {
	preference, offset, err := pack.UnpackUint16(data, offset)
	if err != nil {
		return offset, err
	}
	rr.Preference = preference

	exchange, offset, err := pack.UnpackDomainName(data, offset)
	if err != nil {
		return offset, err
	}
	rr.Exchange = exchange

	return offset, nil
}

func (rr *MX) Pack(buf []byte, offset int, comp compression.Map) (int, error) {
	offset, err := pack.PackUint16(rr.Preference, buf, offset)
	if err != nil {
		return offset, err
	}

	return pack.PackDomainName(rr.Exchange, buf, offset, comp)
}

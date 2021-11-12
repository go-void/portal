package rr

import "github.com/go-void/portal/pkg/wire"

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
	return ""
}

func (rr *MX) Len() uint16 {
	return uint16(len(rr.Exchange)) + 3
}

func (rr *MX) Unpack(data []byte, offset int) (int, error) {
	preference, offset := wire.UnpackUint16(data, offset)
	rr.Preference = preference

	exchange, offset := wire.UnpackDomainName(data, offset)
	rr.Exchange = exchange

	return offset, nil
}

func (rr *MX) Pack(buf []byte, offset int) (int, error) {
	offset, err := wire.PackUint16(rr.Preference, buf, offset)
	if err != nil {
		return offset, err
	}

	return wire.PackDomainName(rr.Exchange, buf, offset)
}
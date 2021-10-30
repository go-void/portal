package rr

import "github.com/go-void/portal/internal/wire"

// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.3.13
type SOA struct {
	H       Header
	MName   string
	RName   string
	Serial  uint32
	Refresh uint32
	Retry   uint32
	Expire  uint32
	Minimum uint32
}

func (rr *SOA) Header() *Header {
	return &rr.H
}

func (rr *SOA) SetHeader(header Header) {
	rr.H = header
}

func (rr *SOA) SetData(data ...interface{}) error {
	if len(data) != 7 {
		return ErrInvalidRRData
	}

	mname, ok := data[0].(string)
	if !ok {
		return ErrFailedToConvertRRData
	}
	rr.MName = mname

	// TODO (Techassi): Add remaining fields

	return nil
}

func (rr *SOA) String() string {
	return ""
}

func (rr *SOA) Unpack(data []byte, offset int) (int, error) {
	mname, offset := wire.UnpackDomainName(data, offset)
	rr.MName = mname

	rname, offset := wire.UnpackDomainName(data, offset)
	rr.RName = rname

	serial, offset := wire.UnpackUint32(data, offset)
	rr.Serial = serial

	refresh, offset := wire.UnpackUint32(data, offset)
	rr.Refresh = refresh

	retry, offset := wire.UnpackUint32(data, offset)
	rr.Retry = retry

	expire, offset := wire.UnpackUint32(data, offset)
	rr.Expire = expire

	minimum, offset := wire.UnpackUint32(data, offset)
	rr.Minimum = minimum

	return offset, nil
}

func (rr *SOA) Pack(buf []byte, offset int) (int, error) {
	offset, err := wire.PackDomainName(rr.MName, buf, offset)
	if err != nil {
		return offset, err
	}

	offset, err = wire.PackDomainName(rr.RName, buf, offset)
	if err != nil {
		return offset, err
	}

	offset, err = wire.PackUint32(rr.Serial, buf, offset)
	if err != nil {
		return offset, err
	}

	offset, err = wire.PackUint32(rr.Refresh, buf, offset)
	if err != nil {
		return offset, err
	}

	offset, err = wire.PackUint32(rr.Retry, buf, offset)
	if err != nil {
		return offset, err
	}

	offset, err = wire.PackUint32(rr.Expire, buf, offset)
	if err != nil {
		return offset, err
	}

	return wire.PackUint32(rr.Minimum, buf, offset)
}

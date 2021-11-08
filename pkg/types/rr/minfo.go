package rr

import "github.com/go-void/portal/pkg/wire"

// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.3.7 (EXPERIMENTAL)
type MINFO struct {
	H        Header
	RMailBox string
	EMailBox string
}

func (rr *MINFO) Header() *Header {
	return &rr.H
}

func (rr *MINFO) SetHeader(header Header) {
	rr.H = header
}

func (rr *MINFO) SetData(data ...interface{}) error {
	if len(data) != 2 {
		return ErrInvalidRRData
	}

	rmail, ok := data[0].(string)
	if !ok {
		return ErrFailedToConvertRRData
	}
	rr.RMailBox = rmail

	email, ok := data[1].(string)
	if !ok {
		return ErrFailedToConvertRRData
	}
	rr.EMailBox = email

	return nil
}

func (rr *MINFO) String() string {
	return ""
}

func (rr *MINFO) Len() uint16 {
	return uint16(len(rr.RMailBox)+len(rr.EMailBox)) + 2
}

func (rr *MINFO) Unpack(data []byte, offset int) (int, error) {
	rmailbox, offset := wire.UnpackDomainName(data, offset)
	rr.RMailBox = rmailbox

	emailbox, offset := wire.UnpackDomainName(data, offset)
	rr.EMailBox = emailbox

	return offset, nil
}

func (rr *MINFO) Pack(buf []byte, offset int) (int, error) {
	offset, err := wire.PackDomainName(rr.RMailBox, buf, offset)
	if err != nil {
		return offset, err
	}

	return wire.PackDomainName(rr.EMailBox, buf, offset)
}

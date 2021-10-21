package rr

// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.3.7 (EXPERIMENTAL)
type MINFO struct {
	H        RRHeader
	RMailBox string
	EMailBox string
}

func (rr *MINFO) Header() *RRHeader {
	return &rr.H
}

func (rr *MINFO) SetHeader(header RRHeader) {
	rr.H = header
}

func (rr *MINFO) SetData(data ...interface{}) error {
	if len(data) != 2 {
		return ErrInvalidRRData
	}

	rmail, ok := data[0].(string)
	if !ok {
		return ErrFailedToConvertRRDate
	}
	rr.RMailBox = rmail

	email, ok := data[1].(string)
	if !ok {
		return ErrFailedToConvertRRDate
	}
	rr.EMailBox = email

	return nil
}

func (rr *MINFO) String() string {
	return ""
}

func (rr *MINFO) Unwrap(data []byte, offset int) (int, error) {
	return offset, nil
}

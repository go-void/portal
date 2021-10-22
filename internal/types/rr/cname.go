package rr

// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.3.1
type CNAME struct {
	H      RRHeader
	Target string
}

func (rr *CNAME) Header() *RRHeader {
	return &rr.H
}

func (rr *CNAME) SetHeader(header RRHeader) {
	rr.H = header
}

func (rr *CNAME) SetData(data ...interface{}) error {
	if len(data) != 1 {
		return ErrInvalidRRData
	}

	target, ok := data[0].(string)
	if !ok {
		return ErrFailedToConvertRRDate
	}
	rr.Target = target
	return nil
}

func (rr *CNAME) String() string {
	return ""
}

func (rr *CNAME) Unpack(data []byte, offset int) (int, error) {
	return offset, nil
}

func (rr *CNAME) Pack(data []byte, offset int) (int, error) {
	return offset, nil
}

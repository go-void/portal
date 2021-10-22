package rr

// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.3.3 (EXPERIMENTAL)
type MB struct {
	H       RRHeader
	MADName string
}

func (rr *MB) Header() *RRHeader {
	return &rr.H
}

func (rr *MB) SetHeader(header RRHeader) {
	rr.H = header
}

func (rr *MB) SetData(data ...interface{}) error {
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

func (rr *MB) String() string {
	return ""
}

func (rr *MB) Unpack(data []byte, offset int) (int, error) {
	return offset, nil
}

func (rr *MB) Pack(data []byte, offset int) (int, error) {
	return offset, nil
}

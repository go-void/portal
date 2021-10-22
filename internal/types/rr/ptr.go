package rr

// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.3.12
type PTR struct {
	H        RRHeader
	PTRDName string
}

func (rr *PTR) Header() *RRHeader {
	return &rr.H
}

func (rr *PTR) SetHeader(header RRHeader) {
	rr.H = header
}

func (rr *PTR) SetData(data ...interface{}) error {
	if len(data) != 1 {
		return ErrInvalidRRData
	}

	name, ok := data[0].(string)
	if !ok {
		return ErrFailedToConvertRRDate
	}
	rr.PTRDName = name
	return nil
}

func (rr *PTR) String() string {
	return ""
}

func (rr *PTR) Unpack(data []byte, offset int) (int, error) {
	return offset, nil
}

func (rr *PTR) Pack(data []byte, offset int) (int, error) {
	return offset, nil
}

package rr

// https://datatracker.ietf.org/doc/html/rfc1035#section-3.3.4 (Obsolete)
type MD struct {
	H       RRHeader
	MADName string
}

func (rr *MD) Header() *RRHeader {
	return &rr.H
}

func (rr *MD) SetHeader(header RRHeader) {
	rr.H = header
}

func (rr *MD) SetData(data ...interface{}) error {
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

func (rr *MD) String() string {
	return ""
}

func (rr *MD) Unwrap(data []byte, offset int) (int, error) {
	return offset, nil
}

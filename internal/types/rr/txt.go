package rr

// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.3.14
type TXT struct {
	H    RRHeader
	Data string
}

func (rr *TXT) Header() *RRHeader {
	return &rr.H
}

func (rr *TXT) SetHeader(header RRHeader) {
	rr.H = header
}

func (rr *TXT) SetData(data ...interface{}) error {
	if len(data) != 1 {
		return ErrInvalidRRData
	}

	d, ok := data[0].(string)
	if !ok {
		return ErrFailedToConvertRRDate
	}
	rr.Data = d
	return nil
}

func (rr *TXT) String() string {
	return ""
}

func (rr *TXT) Unwrap(data []byte, offset int) (int, error) {
	return offset, nil
}

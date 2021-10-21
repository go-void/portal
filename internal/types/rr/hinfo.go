package rr

// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.3.2
type HINFO struct {
	H   RRHeader
	CPU string
	OS  string
}

func (rr *HINFO) Header() *RRHeader {
	return &rr.H
}

func (rr *HINFO) SetHeader(header RRHeader) {
	rr.H = header
}

func (rr *HINFO) SetData(data ...interface{}) error {
	if len(data) != 2 {
		return ErrInvalidRRData
	}

	cpu, ok := data[0].(string)
	if !ok {
		return ErrFailedToConvertRRDate
	}
	rr.CPU = cpu

	os, ok := data[1].(string)
	if !ok {
		return ErrFailedToConvertRRDate
	}
	rr.OS = os

	return nil
}

func (rr *HINFO) String() string {
	return ""
}

func (rr *HINFO) Unwrap(data []byte, offset int) (int, error) {
	return offset, nil
}

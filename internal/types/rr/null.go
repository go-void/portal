package rr

// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.3.10
type NULL struct {
	H    RRHeader
	Data string
}

func (rr *NULL) Header() *RRHeader {
	return &rr.H
}

func (rr *NULL) SetHeader(header RRHeader) {
	rr.H = header
}

func (rr *NULL) SetData(data ...interface{}) error {
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

func (rr *NULL) String() string {
	return ""
}

func (rr *NULL) Unpack(data []byte, offset int) (int, error) {
	return offset, nil
}

func (rr *NULL) Pack(data []byte, offset int) (int, error) {
	// TODO (Techassi): Anything can be put in here. How do we pack that?
	return offset, nil
}

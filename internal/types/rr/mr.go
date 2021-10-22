package rr

// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.3.8 (EXPERIMENTAL)
type MR struct {
	H       RRHeader
	NewName string
}

func (rr *MR) Header() *RRHeader {
	return &rr.H
}

func (rr *MR) SetHeader(header RRHeader) {
	rr.H = header
}

func (rr *MR) SetData(data ...interface{}) error {
	if len(data) != 1 {
		return ErrInvalidRRData
	}

	name, ok := data[0].(string)
	if !ok {
		return ErrFailedToConvertRRDate
	}
	rr.NewName = name

	return nil
}

func (rr *MR) String() string {
	return ""
}

func (rr *MR) Unpack(data []byte, offset int) (int, error) {
	return offset, nil
}

func (rr *MR) Pack(data []byte, offset int) (int, error) {
	return offset, nil
}

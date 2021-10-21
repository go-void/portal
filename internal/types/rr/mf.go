package rr

// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.3.5 (Obsolete)
type MF struct {
	H       RRHeader
	MADName string
}

func (rr *MF) Header() *RRHeader {
	return &rr.H
}

func (rr *MF) SetHeader(header RRHeader) {
	rr.H = header
}

func (rr *MF) SetData(data ...interface{}) error {
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

func (rr *MF) String() string {
	return ""
}

func (rr *MF) Unwrap(data []byte, offset int) (int, error) {
	return offset, nil
}

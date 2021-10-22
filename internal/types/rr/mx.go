package rr

// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.3.9
type MX struct {
	H          RRHeader
	Preference uint16
	Exchange   string
}

func (rr *MX) Header() *RRHeader {
	return &rr.H
}

func (rr *MX) SetHeader(header RRHeader) {
	rr.H = header
}

func (rr *MX) SetData(data ...interface{}) error {
	if len(data) != 2 {
		return ErrInvalidRRData
	}

	pref, ok := data[0].(uint16)
	if !ok {
		return ErrFailedToConvertRRDate
	}
	rr.Preference = pref

	exchange, ok := data[1].(string)
	if !ok {
		return ErrFailedToConvertRRDate
	}
	rr.Exchange = exchange

	return nil
}

func (rr *MX) String() string {
	return ""
}

func (rr *MX) Unpack(data []byte, offset int) (int, error) {
	return offset, nil
}

func (rr *MX) Pack(data []byte, offset int) (int, error) {
	return offset, nil
}

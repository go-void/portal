package rr

// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.3.13
type SOA struct {
	H       RRHeader
	MName   string
	RName   string
	Serial  uint32
	Refresh uint32
	Retry   uint32
	Expire  uint32
	Minimum uint32
}

func (rr *SOA) Header() *RRHeader {
	return &rr.H
}

func (rr *SOA) SetHeader(header RRHeader) {
	rr.H = header
}

func (rr *SOA) SetData(data ...interface{}) error {
	if len(data) != 7 {
		return ErrInvalidRRData
	}

	mname, ok := data[0].(string)
	if !ok {
		return ErrFailedToConvertRRDate
	}
	rr.MName = mname

	// TODO (Techassi): Add remaining fields

	return nil
}

func (rr *SOA) String() string {
	return ""
}

func (rr *SOA) Unpack(data []byte, offset int) (int, error) {
	return offset, nil
}

func (rr *SOA) Pack(data []byte, offset int) (int, error) {
	return offset, nil
}

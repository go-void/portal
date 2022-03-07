package rr

import (
	"errors"

	"github.com/go-void/portal/pkg/constants"
	"github.com/go-void/portal/pkg/pack"
)

var (
	ErrSOASerialOutOfRange = errors.New("serial out of range")
)

type SerialComparison int

const (
	SerialEqual SerialComparison = iota
	SerialLess
	SerialGreater
)

// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.3.13
type SOA struct {
	H       Header
	MName   string
	RName   string
	Serial  uint32
	Refresh uint32
	Retry   uint32
	Expire  uint32
	Minimum uint32
}

func (rr *SOA) Header() *Header {
	return &rr.H
}

func (rr *SOA) SetHeader(header Header) {
	rr.H = header
}

func (rr *SOA) SetData(data ...interface{}) error {
	if len(data) != 7 {
		return ErrInvalidRRData
	}

	mname, ok := data[0].(string)
	if !ok {
		return ErrFailedToConvertRRData
	}
	rr.MName = mname

	rname, ok := data[1].(string)
	if !ok {
		return ErrFailedToConvertRRData
	}
	rr.RName = rname

	serial, ok := data[2].(uint32)
	if !ok {
		return ErrFailedToConvertRRData
	}
	rr.Serial = serial

	refresh, ok := data[3].(uint32)
	if !ok {
		return ErrFailedToConvertRRData
	}
	rr.Refresh = refresh

	retry, ok := data[4].(uint32)
	if !ok {
		return ErrFailedToConvertRRData
	}
	rr.Retry = retry

	expire, ok := data[5].(uint32)
	if !ok {
		return ErrFailedToConvertRRData
	}
	rr.Expire = expire

	minimum, ok := data[6].(uint32)
	if !ok {
		return ErrFailedToConvertRRData
	}
	rr.Minimum = minimum

	return nil
}

func (rr *SOA) String() string {
	return ""
}

func (rr *SOA) Len() uint16 {
	return uint16(len(rr.MName)+len(rr.RName)) + 2 + 20
}

func (rr *SOA) IsSame(o RR) bool {
	other, ok := o.(*SOA)
	if !ok {
		return false
	}

	if rr.MName != other.MName {
		return false
	}
	if rr.RName != other.RName {
		return false
	}
	if rr.Serial != other.Serial {
		return false
	}
	if rr.Refresh != other.Refresh {
		return false
	}
	if rr.Retry != other.Retry {
		return false
	}
	if rr.Expire != other.Expire {
		return false
	}

	return rr.Minimum == other.Minimum
}

func (rr *SOA) Unpack(data []byte, offset int) (int, error) {
	mname, offset, err := pack.UnpackDomainName(data, offset)
	if err != nil {
		return offset, err
	}
	rr.MName = mname

	rname, offset, err := pack.UnpackDomainName(data, offset)
	if err != nil {
		return offset, err
	}
	rr.RName = rname

	serial, offset, err := pack.UnpackUint32(data, offset)
	if err != nil {
		return offset, err
	}
	rr.Serial = serial

	refresh, offset, err := pack.UnpackUint32(data, offset)
	if err != nil {
		return offset, err
	}
	rr.Refresh = refresh

	retry, offset, err := pack.UnpackUint32(data, offset)
	if err != nil {
		return offset, err
	}
	rr.Retry = retry

	expire, offset, err := pack.UnpackUint32(data, offset)
	if err != nil {
		return offset, err
	}
	rr.Expire = expire

	minimum, offset, err := pack.UnpackUint32(data, offset)
	if err != nil {
		return offset, err
	}
	rr.Minimum = minimum

	return offset, nil
}

func (rr *SOA) Pack(buf []byte, offset int) (int, error) {
	offset, err := pack.PackDomainName(rr.MName, buf, offset)
	if err != nil {
		return offset, err
	}

	offset, err = pack.PackDomainName(rr.RName, buf, offset)
	if err != nil {
		return offset, err
	}

	offset, err = pack.PackUint32(rr.Serial, buf, offset)
	if err != nil {
		return offset, err
	}

	offset, err = pack.PackUint32(rr.Refresh, buf, offset)
	if err != nil {
		return offset, err
	}

	offset, err = pack.PackUint32(rr.Retry, buf, offset)
	if err != nil {
		return offset, err
	}

	offset, err = pack.PackUint32(rr.Expire, buf, offset)
	if err != nil {
		return offset, err
	}

	return pack.PackUint32(rr.Minimum, buf, offset)
}

func (rr *SOA) SerialAdd(n int) error {
	if n < 0 || n > constants.SerialMaxAdditon {
		return ErrSOASerialOutOfRange
	}

	rr.Serial = (rr.Serial + uint32(n)) % constants.SerialBits
	return nil
}

func (rr *SOA) SerialCompare(soa *SOA) SerialComparison {
	if rr.Serial == soa.Serial {
		return SerialEqual
	}

	if (rr.Serial < soa.Serial && soa.Serial-rr.Serial < constants.SerialMaxBits) ||
		(rr.Serial > soa.Serial && rr.Serial-soa.Serial > constants.SerialMaxBits) {
		return SerialLess
	}

	if (rr.Serial < soa.Serial && soa.Serial-rr.Serial > constants.SerialMaxBits) ||
		(rr.Serial > soa.Serial && rr.Serial-soa.Serial < constants.SerialMaxBits) {
		return SerialGreater
	}

	return SerialEqual
}

package edns

import (
	"fmt"
)

// DNSSEC Algorithm Understood
// See https://datatracker.ietf.org/doc/html/rfc6975#section-3
type DAU struct {
	Algorithms []uint16
}

// Code returns the option code
func (o *DAU) Code() uint16 {
	return CodeDAU
}

// Len returns the option length
func (o *DAU) Len() uint16 {
	return 0
}

func (o *DAU) String() string {
	return fmt.Sprintf("DAU: %v", o.Algorithms)
}

// Unpack unpacks the option data
func (o *DAU) Unpack(data []byte, offset int, length uint16) (int, error) {
	return offset, nil
}

// Pack packs the option data
func (o *DAU) Pack(buf []byte, offset int) (int, error) {
	return offset, nil
}

package edns

import (
	"fmt"
)

// NSEC3 Hash Understood
// See https://datatracker.ietf.org/doc/html/rfc6975#section-3
type N3U struct {
	Algorithms []uint16
}

// Code returns the option code
func (o *N3U) Code() uint16 {
	return CodeN3U
}

// Len returns the option length
func (o *N3U) Len() uint16 {
	return 0
}

func (o *N3U) String() string {
	return fmt.Sprintf("DHU: %v", o.Algorithms)
}

// Unpack unpacks the option data
func (o *N3U) Unpack(data []byte, offset int, length uint16) (int, error) {
	return offset, nil
}

// Pack packs the option data
func (o *N3U) Pack(buf []byte, offset int) (int, error) {
	return offset, nil
}

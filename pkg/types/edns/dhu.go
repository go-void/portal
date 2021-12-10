package edns

import (
	"fmt"
)

// DS Hash Understood
// See https://datatracker.ietf.org/doc/html/rfc6975#section-3
type DHU struct {
	Algorithms []uint16
}

// Code returns the option code
func (o *DHU) Code() uint16 {
	return CodeDHU
}

// Len returns the option length
func (o *DHU) Len() uint16 {
	return 0
}

func (o *DHU) String() string {
	return fmt.Sprintf("DHU: %v", o.Algorithms)
}

// Unpack unpacks the option data
func (o *DHU) Unpack(data []byte, offset int, length uint16) (int, error) {
	return offset, nil
}

// Pack packs the option data
func (o *DHU) Pack(buf []byte, offset int) (int, error) {
	return offset, nil
}

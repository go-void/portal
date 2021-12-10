package edns

import (
	"encoding/hex"
	"fmt"
)

type NSID struct {
	ID string
}

// Code returns the option code
func (o *NSID) Code() uint16 {
	return CodeNSID
}

// Len returns the option length
func (o *NSID) Len() uint16 {
	return uint16(len(o.ID))
}

func (o *NSID) String() string {
	return fmt.Sprintf("NSID: %s", o.ID)
}

// Unpack unpacks the option data
func (o *NSID) Unpack(data []byte, offset int, length uint16) (int, error) {
	id := data[offset : offset+int(length)]
	o.ID = hex.EncodeToString(id)
	return offset + int(length), nil
}

// Pack packs the option data
func (o *NSID) Pack(buf []byte, offset int) (int, error) {
	id, err := hex.DecodeString(o.ID)
	if err != nil {
		return offset, err
	}
	copy(buf[offset:], id)
	return offset + len(id), nil
}

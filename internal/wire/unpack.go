package wire

import (
	"encoding/binary"
)

// TODO (Techassi): Add offset overflow check

// UnpackUint8 unpacks a uint8 of 'data' at 'offset'
// and returns the new offset
func UnpackUint8(data []byte, offset int) (uint8, int) {
	return data[offset], offset + 1
}

// UnpackUint16 unpacks a uint16 of 'data' at 'offset'
// and returns the new offset
func UnpackUint16(data []byte, offset int) (uint16, int) {
	return binary.BigEndian.Uint16(data[offset:]), offset + 2
}

// UnpackUint32 unpacks a uint32 of 'data' at 'offset'
// and returns the new offset
func UnpackUint32(data []byte, offset int) (uint32, int) {
	return binary.BigEndian.Uint32(data[offset:]), offset + 4
}

// FIXME (Techassi): Add ability to follow pointers marked by upper 2 bits of an octet set to 1 followed by 14 bits of
// offset (6 bits from first octet + one octet)

// UnpackDomainName unwraps a domain name in a DNS question or in a RR header
func UnpackDomainName(data []byte, offset int) (string, int) {
	// If we immediation encounter a null byte, the name is root (.)
	if data[offset] == 0x00 {
		return ".", offset + 1
	}

	var offsetBeforePtr = 0
	var followed bool
	var buf []byte
	var done bool

	for !done {
		b := int(data[offset])
		offset++

		// Check if we have a pointer (11000000 => 0xC0)
		switch b & 0xC0 {
		case 0x00:
			if b == 0x00 {
				done = true
				break
			}

			buf = append(buf, data[offset:offset+b]...)
			buf = append(buf, '.')
			offset += b
		case 0xC0:
			if !followed {
				offsetBeforePtr = offset + 1
			}
			offset = int(b^0xC0)<<8 | int(data[offset])
			followed = true
		}
	}

	if followed {
		offset = offsetBeforePtr
	}

	return string(buf), offset
}

func UnpackCharacterString(data []byte, offset int) (int, error) {
	return offset, nil
}

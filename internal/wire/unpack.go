package wire

import "encoding/binary"

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

// UnpackDomainName unwraps a domain name in a DNS question or in a RR header
func UnpackDomainName(data []byte, offset int) (string, int) {
	// If we immediation encounter a null byte, the name is root (.)
	if data[offset] == 0x00 {
		return ".", offset + 1
	}

	var buf []byte

	// Initialize the end of the label in bytes
	end := offset + int(data[offset]) + 1
	offset++

	// Iterate over the bytes until we reach the null byte, which
	// marks the root (.)
	for {
		if data[offset] == 0x00 {
			buf = append(buf, '.')
			offset++
			break
		}

		if offset == end {
			buf = append(buf, '.')
			end += int(data[offset]) + 1
			offset++
		}

		buf = append(buf, data[offset])
		offset++
	}

	return string(buf), offset
}

func UnpackCharacterString(data []byte, offset int) (int, error) {
	return offset, nil
}

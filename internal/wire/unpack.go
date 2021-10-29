package wire

import (
	"encoding/binary"
	"net"
)

// TODO (Techassi): Add offset overflow check

// UnpackUint8 unpacks a uint8 from data at offset and returns the new offset
func UnpackUint8(data []byte, offset int) (uint8, int) {
	return data[offset], offset + 1
}

// UnpackUint16 unpacks a uint16 from data at offset and returns the new offset
func UnpackUint16(data []byte, offset int) (uint16, int) {
	return binary.BigEndian.Uint16(data[offset:]), offset + 2
}

// UnpackUint32 unpacks a uint32 from data at offset and returns the new offset
func UnpackUint32(data []byte, offset int) (uint32, int) {
	return binary.BigEndian.Uint32(data[offset:]), offset + 4
}

// UnpackUint64 unpacks a uint64 from data at offset and returns the new offset
func UnpackUint64(data []byte, offset int) (uint64, int) {
	return binary.BigEndian.Uint64(data[offset:]), offset + 8
}

// UnpackIPv4Address unpacks a IPv4 address and returns the new offset
func UnpackIPv4Address(data []byte, offset int) (net.IP, int) {
	return net.IPv4(data[offset], data[offset+1], data[offset+2], data[offset+3]), offset + 4
}

// NOTE (Techassi): Is this actually correct?

// UnpackIPv6Address unpacks a IPv6 address and returns the new offset
func UnpackIPv6Address(data []byte, offset int) (net.IP, int) {
	hi, offset := UnpackUint64(data, offset)
	lo, offset := UnpackUint64(data, offset)

	ip := make(net.IP, net.IPv6len)
	binary.BigEndian.PutUint64(ip, hi)
	binary.BigEndian.PutUint64(ip[8:], lo)

	return ip, offset
}

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

package pack

import (
	"encoding/binary"
	"net"

	"github.com/go-void/portal/pkg/types/edns"
)

// UnpackUint8 unpacks a uint8 from data at offset and returns the new offset
func UnpackUint8(data []byte, offset int) (uint8, int, error) {
	if offset+1 > len(data) {
		return 0, len(data), ErrOverflowUnpackUint8
	}
	return data[offset], offset + 1, nil
}

// UnpackUint16 unpacks a uint16 from data at offset and returns the new offset
func UnpackUint16(data []byte, offset int) (uint16, int, error) {
	if offset+2 > len(data) {
		return 0, len(data), ErrOverflowUnpackUint16
	}
	return binary.BigEndian.Uint16(data[offset:]), offset + 2, nil
}

// UnpackUint32 unpacks a uint32 from data at offset and returns the new offset
func UnpackUint32(data []byte, offset int) (uint32, int, error) {
	if offset+4 > len(data) {
		return 0, len(data), ErrOverflowUnpackUint32
	}
	return binary.BigEndian.Uint32(data[offset:]), offset + 4, nil
}

// UnpackUint64 unpacks a uint64 from data at offset and returns the new offset
func UnpackUint64(data []byte, offset int) (uint64, int, error) {
	if offset+8 > len(data) {
		return 0, len(data), ErrOverflowUnpackUint64
	}
	return binary.BigEndian.Uint64(data[offset:]), offset + 8, nil
}

// UnpackIPv4Address unpacks a IPv4 address and returns the new offset
func UnpackIPv4Address(data []byte, offset int) (net.IP, int, error) {
	if offset+4 > len(data) {
		return net.IP{}, len(data), ErrOverflowUnpackIPv4
	}
	return net.IPv4(data[offset], data[offset+1], data[offset+2], data[offset+3]), offset + 4, nil
}

// UnpackIPv6Address unpacks a IPv6 address and returns the new offset
func UnpackIPv6Address(data []byte, offset int) (net.IP, int, error) {
	hi, offset, err := UnpackUint64(data, offset)
	if err != nil {
		return net.IP{}, len(data), ErrOverflowUnpackIPv6
	}

	lo, offset, err := UnpackUint64(data, offset)
	if err != nil {
		return net.IP{}, len(data), ErrOverflowUnpackIPv6
	}

	ip := make(net.IP, net.IPv6len)
	binary.BigEndian.PutUint64(ip, hi)
	binary.BigEndian.PutUint64(ip[8:], lo)

	return ip, offset, nil
}

// UnpackDomainName unoacks a domain name in a DNS question or in a RR header
func UnpackDomainName(data []byte, offset int) (string, int, error) {
	dataLength := len(data)
	if offset > dataLength {
		return "", dataLength, ErrOverflowUnpackName
	}

	// If we immediatly encounter a null byte, the name is root (.)
	if data[offset] == 0x00 {
		return ".", offset + 1, nil
	}

	var offsetBeforePtr = 0
	var followed bool
	var buf []byte

	for done := false; !done; {
		b := int(data[offset])
		offset++

		// Check if we have a pointer (11000000 => 0xC0). Pointers point
		// to domain names previously defined in some part of the message.
		// We follow the pointer (by updating the offset) and reading in
		// the domain name as usual. After encountering the null byte we
		// jump back by updating the offset
		switch b & 0xC0 {
		case 0x00:
			if b == 0x00 {
				done = true
				break
			}

			if offset+b > dataLength {
				return "", dataLength, ErrOverflowUnpackName
			}

			buf = append(buf, data[offset:offset+b]...)
			buf = append(buf, '.')
			offset += b
		case 0xC0:
			if offset+b > dataLength {
				return "", dataLength, ErrOverflowUnpackName
			}

			if !followed {
				offsetBeforePtr = offset + 1
			}

			offset = (b^0xC0)<<8 | int(data[offset])
			followed = true
		}
	}

	if followed {
		offset = offsetBeforePtr
	}

	return string(buf), offset, nil
}

// UnpackCharacterString unpacks a character string.
// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.3 <character-string>
func UnpackCharacterString(data []byte, offset int) (string, int, error) {
	l := int(data[offset])

	if offset+l > len(data) {
		return "", len(data), ErrOverflowUnpackString
	}

	t := make([]byte, l)
	offset++

	for i := 0; i < l; i++ {
		t[i] = data[offset]
		offset++
	}

	return string(t), offset, nil
}

func UnpackEDNSOptions(data []byte, offset int, rdlen uint16) ([]edns.Option, int, error) {
	var optionLength = offset + int(rdlen)
	var options []edns.Option

	for offset < optionLength {
		// TODO (Techassi): Add offset validation
		code := binary.BigEndian.Uint16(data[offset:])
		offset += 2

		length := binary.BigEndian.Uint16(data[offset:])
		offset += 2

		option, err := edns.New(code)
		if err != nil {
			return options, offset, err
		}

		// TODO (Techassi): Add offset validation
		offset, err = option.Unpack(data, offset, length)
		if err != nil {
			return options, offset, err
		}

		options = append(options, option)
		offset += int(length)
	}

	return options, offset, nil
}

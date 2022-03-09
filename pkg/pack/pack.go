package pack

import (
	"encoding/binary"
	"errors"
	"net"

	"github.com/go-void/portal/pkg/compression"
	"github.com/go-void/portal/pkg/types/edns"
)

var (
	ErrCharacterStringTooLong = errors.New("pack: character string too long")
	ErrInvalidName            = errors.New("pack: invalid name")
)

// PackUint8 packs a uint8 (one octet) into buf and returns the new offset
func PackUint8(data uint8, buf []byte, offset int) (int, error) {
	if offset+1 > len(buf) {
		return len(buf), ErrOverflowPackUint8
	}

	buf[offset] = data
	return offset + 1, nil
}

// PackUint16 packs a uint16 (two octets) into buf and returns the new offset
func PackUint16(data uint16, buf []byte, offset int) (int, error) {
	if offset+2 > len(buf) {
		return len(buf), ErrOverflowPackUint16
	}

	binary.BigEndian.PutUint16(buf[offset:], data)
	return offset + 2, nil
}

// PackUint32 ppacks a uint32 (four octets) into buff and returns the new offset
func PackUint32(data uint32, buf []byte, offset int) (int, error) {
	if offset+4 > len(buf) {
		return len(buf), ErrOverflowPackUint32
	}

	binary.BigEndian.PutUint32(buf[offset:], data)
	return offset + 4, nil
}

// PackIPAddress packs a IP address (v4 or v6) into buf and returns the new offset
func PackIPAddress(ip net.IP, buf []byte, offset int) (int, error) {
	if i := ip.To4(); i != nil {
		if offset+net.IPv4len > len(buf) {
			return len(buf), ErrOverflowPackIPv4
		}

		ipAsUint32 := binary.BigEndian.Uint32(i)
		binary.BigEndian.PutUint32(buf[offset:], ipAsUint32)
		return offset + 4, nil
	}

	if offset+net.IPv6len > len(buf) {
		return len(buf), ErrOverfloPpackIPv6
	}

	i := binary.BigEndian.Uint64(ip[:8])
	binary.BigEndian.PutUint64(buf[offset:], i)
	offset += 8

	i = binary.BigEndian.Uint64(ip[8:])
	binary.BigEndian.PutUint64(buf[offset:], i)

	return offset + 8, nil
}

// TODO (Techassi): Add offset overflow check
// TODO (Techassi): Implement message compression

// PackDomainName packs a name into buf and returns the new offset.
func PackDomainName(name string, buf []byte, offset int, comp compression.Map) (int, error) {
	length := len(name)
	dot := false
	pos := 0

	for i := 0; i < length; i++ {
		b := name[i]
		switch {
		case b == '.':
			// Two dots after each other are invalid, return
			if dot {
				return len(buf), ErrInvalidName
			}
			dot = true

			// If name is root, break
			if i == 0 {
				break
			}

			// TODO (Techassi): Add overflow check

			// TODO (Techassi): Handle compression
			comp.Set(name[pos:i], offset)

			// Append the label length to the buffer
			labelLength := i - pos
			buf[offset] = byte(labelLength)
			offset++

			// Add individual bytes of the label to the buffer
			copy(buf[offset:], name[pos:i])

			// Update offset and position
			offset += labelLength
			pos = i + 1
		case b == '-',
			b >= 0x30 && b <= 0x39, // ASCII 0-9
			b >= 0x41 && b <= 0x5A, // ASCII A-Z
			b >= 0x61 && b <= 0x7A: // ASCII a-z
			dot = false
		default:
			return len(buf), ErrInvalidName
		}
	}

	// We packed the complete name, add null byte
	buf[offset] = 0x0
	return offset + 1, nil
}

// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.3 <character-string>
func PackCharacterString(characters string, buf []byte, offset int) (int, error) {
	// Character strings can only be 256 octets long (including the length octet)
	if len(characters) > 255 {
		return offset, ErrCharacterStringTooLong
	}

	buf[offset] = uint8(len(characters))
	offset++

	for i := 0; i < len(characters); i++ {
		buf[offset] = characters[i]
		offset++
	}

	return offset, nil
}

// PackEDNSOptions packs all EDNS options into buf and returns the new offset
func PackEDNSOptions(options []edns.Option, buf []byte, offset int) (int, error) {
	for _, option := range options {
		o, err := PackUint16(option.Code(), buf, offset)
		if err != nil {
			return offset, err
		}
		offset = o

		o, err = PackUint16(option.Len(), buf, offset)
		if err != nil {
			return offset, err
		}
		offset = o

		o, err = option.Pack(buf, offset)
		if err != nil {
			return offset, err
		}
		offset = o
	}

	return offset, nil
}

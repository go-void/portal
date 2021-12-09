package pack

import (
	"encoding/binary"
	"errors"
	"net"

	"github.com/go-void/portal/pkg/labels"
	"github.com/go-void/portal/pkg/types/edns"
)

var (
	ErrCharacterStringTooLong = errors.New("character string too long")
)

// TODO (Techassi): Add offset overflow check

// PackUint8 packs a uint8 (one octet) into buf and returns the
// new offset
func PackUint8(data uint8, buf []byte, offset int) (int, error) {
	buf[offset] = data
	return offset + 1, nil
}

// PackUint16 packs a uint16 (two octets) into buf and returns the
// new offset
func PackUint16(data uint16, buf []byte, offset int) (int, error) {
	binary.BigEndian.PutUint16(buf[offset:], data)
	return offset + 2, nil
}

// PackUint32 ppacks a uint32 (four octets) into buff and returns the
// new offset
func PackUint32(data uint32, buf []byte, offset int) (int, error) {
	binary.BigEndian.PutUint32(buf[offset:], data)
	return offset + 4, nil
}

// PackIPAddress packs a IP address (v4 or v6) into buf and returns
// the new offset
func PackIPAddress(ip net.IP, buf []byte, offset int) (int, error) {
	if len(ip) == 16 {
		i := binary.BigEndian.Uint32(ip[12:16])
		binary.BigEndian.PutUint32(buf[offset:], i)
		return offset + 4, nil
	}

	i := binary.BigEndian.Uint64(ip[:32])
	binary.BigEndian.PutUint64(buf[offset:], i)
	offset += 8

	i = binary.BigEndian.Uint64(ip[32:])
	binary.BigEndian.PutUint64(buf[offset:], i)

	return offset + 8, nil
}

// TODO (Techassi): Implement message compression

// PackDomainName packs a name into buf and returns the new offset.
// Example: example.com => 7 101 120 97 109 112 108 101 3 99 111 109 0.
// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.3 <domain-name>
func PackDomainName(name string, buf []byte, offset int) (int, error) {
	labels := labels.FromBottom(name)

	for i := 0; i < len(labels); i++ {
		label := labels[i]
		switch label {
		case "", ".":
			break
		default:
			buf[offset] = uint8(len(label))
			offset++

			for l := 0; l < len(label); l++ {
				buf[offset] = label[l]
				offset++
			}
			continue
		}

		buf[offset] = 0x0
		offset++
	}

	return offset, nil
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

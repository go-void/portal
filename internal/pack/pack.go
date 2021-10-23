package pack

import (
	"encoding/binary"
	"net"
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

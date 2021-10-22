package utils

import "encoding/binary"

// TODO (Techassi): Add offset overflow check

// PackUint8 unwraps a uint8 of 'data' at 'offset' and returns the
// new offset
func PackUint8(data uint8, buf []byte, offset int) (int, error) {
	buf[offset] = data
	return offset + 1, nil
}

// PackUint16 unwraps a uint16 of 'data' at 'offset' and returns the
// new offset
func PackUint16(data uint16, buf []byte, offset int) (int, error) {
	binary.BigEndian.PutUint16(buf[offset:], data)
	return offset + 2, nil
}

// PackUint32 unwraps a uint32 of 'data' at 'offset' and returns the
// new offset
func PackUint32(data uint32, buf []byte, offset int) (int, error) {
	binary.BigEndian.PutUint32(buf[offset:], data)
	return offset + 4, nil
}

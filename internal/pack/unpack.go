package pack

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

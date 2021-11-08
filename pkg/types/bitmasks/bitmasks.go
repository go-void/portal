package bitmasks

// This file defines bitmasks which get used to extract
// DNS header flags from octets

const (
	QR = 1 << 15
	AA = 1 << 10
	TC = 1 << 9
	RD = 1 << 8
	RA = 1 << 7
	Z  = 1 << 6
)

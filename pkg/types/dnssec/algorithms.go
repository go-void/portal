package dnssec

// See https://www.iana.org/assignments/dns-sec-alg-numbers/dns-sec-alg-numbers.xhtml

const (
	AlgoDELETE uint16 = 0
	AlgoRSAMD5 uint16 = 1 // Deprecated
	AlgoDH     uint16 = 2 // Diffie-Hellman
	AlgoDSA    uint16 = 4 // DSA/SHA1

	// 4 is reserved

	AlgoRSASHA1          uint16 = 5 // RSA/SHA-1
	AlgoDSANSEC3SHA1     uint16 = 6 // DSA-NSEC3-SHA1
	AlgoRSASHA1NSEC3SHA1 uint16 = 7 // RSASHA1-NSEC3-SHA1
	AlgoRSASHA256        uint16 = 8 // RSASHA256

	// 9 is reserved

	AlgoRSASHA512 uint16 = 10 // RSA/SHA-512

	// 11 is reserved

	AlgoECC_GOST        uint16 = 12 // GOST R 34.10-2001
	AlgoECDSAP256SHA256 uint16 = 13 // ECDSA Curve P-256 with SHA-256
	AlgoECDSAP384SHA384 uint16 = 14 // ECDSA Curve P-384 with SHA-384
	AlgoED25519         uint16 = 15 // Ed25519
	AlgoED448           uint16 = 16 // Ed448

	// 17 - 122 are unassigned
	// 123 - 251 are reserved

	AlgoINDIRECT   uint16 = 252 // Reserved for Indirect Keys
	AlgoPRIVATEDNS uint16 = 253 // private algorithm
	AlgoPRIVATEOID uint16 = 254 // private algorithm OID

	// 255 is reserved
)

package constants

const (
	// DNS questions always have a minimum fixed length of 2 octects for QTYPE and QCLASS
	DNSQuestionFixedLen = 4

	// The DNS header is always 12 octects long
	DNSHeaderLen = 12
)

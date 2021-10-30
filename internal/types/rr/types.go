package rr

// Type describes RR type codes.
// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.2.2
type Type uint16

const (
	TypeNone  uint16 = 0
	TypeA     uint16 = 1  // A host address
	TypeNS    uint16 = 2  // An authoritative name server
	TypeMD    uint16 = 3  // A mail destination (Obsolete - use MX)
	TypeMF    uint16 = 4  // A mail forwarded (Obsolete - use MX)
	TypeCNAME uint16 = 5  // The canonical name for an alias
	TypeSOA   uint16 = 6  // Marks the start of a zone of authority
	TypeMB    uint16 = 7  // A mailbox domain name (EXPERIMENTAL)
	TypeMG    uint16 = 8  // A mail group member (EXPERIMENTAL)
	TypeMR    uint16 = 9  // A mail rename domain name (EXPERIMENTAL)
	TypeNULL  uint16 = 10 // A null RR (EXPERIMENTAL)
	TypePTR   uint16 = 12 // A domain name pointer
	TypeHINFO uint16 = 13 // Host information
	TypeMINFO uint16 = 14 // Mailbox or mail list information
	TypeMX    uint16 = 15 // Mail exchange
	TypeTXT   uint16 = 16 // Text strings

	TypeAAAA uint16 = 28 // AAAA host address

	// QTypes are a superset of types and should only be
	// allowed in questions

	TypeAXFR  uint16 = 252 // A request for a transfer of an entire zone
	TypeMAILB uint16 = 253 // A request for mailbox-related records (MB, MG or MR)
	TypeMAILA uint16 = 254 // A request for mail agent RRs (Obsolete - see MX)
	TypeANY   uint16 = 255 // A request for all records (*)
)

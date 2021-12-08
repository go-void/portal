package rr

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
	TypeAAAA  uint16 = 28 // AAAA host address
	TypeOPT   uint16 = 41 // OPT Record / Meta record

	// QTypes are a superset of types and should only be
	// allowed in questions

	TypeAXFR  uint16 = 252 // A request for a transfer of an entire zone
	TypeMAILB uint16 = 253 // A request for mailbox-related records (MB, MG or MR)
	TypeMAILA uint16 = 254 // A request for mail agent RRs (Obsolete - see MX)
	TypeANY   uint16 = 255 // A request for all records (*)
)

var typeMap = map[uint16]func() RR{
	TypeCNAME: func() RR { return new(CNAME) },
	TypeHINFO: func() RR { return new(HINFO) },
	TypeMB:    func() RR { return new(MB) },
	TypeMD:    func() RR { return new(MD) },
	TypeMF:    func() RR { return new(MF) },
	TypeMG:    func() RR { return new(MG) },
	TypeMINFO: func() RR { return new(MINFO) },
	TypeMR:    func() RR { return new(MR) },
	TypeMX:    func() RR { return new(MX) },
	TypeNULL:  func() RR { return new(NULL) },
	TypeNS:    func() RR { return new(NS) },
	TypePTR:   func() RR { return new(PTR) },
	TypeSOA:   func() RR { return new(SOA) },
	TypeTXT:   func() RR { return new(TXT) },
	TypeA:     func() RR { return new(A) },
	TypeAAAA:  func() RR { return new(AAAA) },
}

var typeToStringMap = map[uint16]string{
	TypeNone:  "NONE",
	TypeA:     "A",
	TypeNS:    "NS",
	TypeMD:    "MD",
	TypeMF:    "MF",
	TypeCNAME: "CNAME",
	TypeSOA:   "SOA",
	TypeMB:    "MB",
	TypeMG:    "MG",
	TypeMR:    "MR",
	TypeNULL:  "NULL",
	TypePTR:   "PTR",
	TypeHINFO: "HINFO",
	TypeMINFO: "MINFO",
	TypeMX:    "MX",
	TypeTXT:   "TXT",
	TypeAAAA:  "AAAA",
	TypeOPT:   "OPT",
	TypeAXFR:  "AXFR",
	TypeMAILB: "MAILB",
	TypeMAILA: "MAILA",
	TypeANY:   "ANY",
}

var stringToTypeMap = map[string]uint16{
	"NONE":  TypeNone,
	"A":     TypeA,
	"NS":    TypeNS,
	"MD":    TypeMD,
	"MF":    TypeMF,
	"CNAME": TypeCNAME,
	"SOA":   TypeSOA,
	"MB":    TypeMB,
	"MG":    TypeMG,
	"MR":    TypeMR,
	"NULL":  TypeNULL,
	"PTR":   TypePTR,
	"HINFO": TypeHINFO,
	"MINFO": TypeMINFO,
	"MX":    TypeMX,
	"TXT":   TypeTXT,
	"AAAA":  TypeAAAA,
	"OPT":   TypeOPT,
	"AXFR":  TypeAXFR,
	"MAILB": TypeMAILB,
	"MAILA": TypeMAILA,
	"ANY":   TypeANY,
}

func TypeToString(t uint16) string {
	return typeToStringMap[t]
}

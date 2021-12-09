package edns

const (
	// 0 is reserved

	CodeLLQ  uint16 = 1 // Long-Lived Queries [RFC 8764]
	CodeUL   uint16 = 2 // Update Leases (Draft)
	CodeNSID uint16 = 3 // Name Server Identifier [RFC 5001]

	// 4 is reserved

	CodeDAU          uint16 = 5
	CodeDHU          uint16 = 6
	CodeN3U          uint16 = 7
	CodeECS          uint16 = 8
	CodeEXPIRE       uint16 = 9
	CodeCOOKIE       uint16 = 10 // Cookies [RFC 7873]
	CodeTCPKEEPALIVE uint16 = 11
	CodePADDING      uint16 = 12
	CodeCHAIN        uint16 = 13
	CodeKEYTAG       uint16 = 14
	CodeEDE          uint16 = 15
	CodeCLIENTTAG    uint16 = 16
	CodeSERVERTAG    uint16 = 17

	// 18-20291 unassigned

	CodeUMBRELLAIDENT uint16 = 20292

	// 20293-26945 unassigned

	CodeDEVICEID uint16 = 26946

	// 26947-65000 unassigned
	// 65001-65534 reserved for local/experimental use
	// 65535 reserved for future expansion
)

var codeMap = map[uint16]func() Option{
	CodeLLQ:           func() Option { return nil },
	CodeUL:            func() Option { return nil },
	CodeNSID:          func() Option { return nil },
	CodeDAU:           func() Option { return nil },
	CodeDHU:           func() Option { return nil },
	CodeN3U:           func() Option { return nil },
	CodeECS:           func() Option { return nil },
	CodeEXPIRE:        func() Option { return nil },
	CodeCOOKIE:        func() Option { return new(Cookie) },
	CodeTCPKEEPALIVE:  func() Option { return nil },
	CodePADDING:       func() Option { return nil },
	CodeCHAIN:         func() Option { return nil },
	CodeKEYTAG:        func() Option { return nil },
	CodeEDE:           func() Option { return nil },
	CodeCLIENTTAG:     func() Option { return nil },
	CodeSERVERTAG:     func() Option { return nil },
	CodeUMBRELLAIDENT: func() Option { return nil },
	CodeDEVICEID:      func() Option { return nil },
}

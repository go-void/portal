package rr

// Class describes RR class codes.
// See https://datatracker.ietf.org/doc/html/rfc1035#section-3.2.4
const (
	IN uint16 = 1 // The Internet
	CS uint16 = 2 // The CSNET class (Obsolete - used only for examples in some obsolete RFCs)
	CH uint16 = 3 // The CHAOS class
	HS uint16 = 4 // Hesiod [Dyer 87]

	// QClasses are a superset of classes and should only be
	// allowed in questions

	ANY uint16 = 255 // Any class (*)
)

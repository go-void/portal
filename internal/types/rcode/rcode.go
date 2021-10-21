package rcode

// RCode describes the kind of response
type Code uint16

const (
	NoError Code = iota
	FormatError
	ServerFailure
	NameError
	NotImplemented
	Refused
)

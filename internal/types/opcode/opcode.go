package opcode

// OpCode describes the kind of query of the message
type Code uint16

const (
	Query Code = iota
	IQuery
	Status
)

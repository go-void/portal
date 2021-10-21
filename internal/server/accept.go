package server

import (
	"github.com/go-void/portal/internal/types/dns"
	"github.com/go-void/portal/internal/types/opcode"
)

type AcceptFunc func(dns.MessageHeader) AcceptAction

type AcceptAction int

const (
	AcceptMessage AcceptAction = iota
	RejectMessage
	IgnoreMessage
	NoImplMessage
)

func DefaultAcceptFunc(h dns.MessageHeader) AcceptAction {
	if !h.IsQuery {
		return IgnoreMessage
	}

	if h.OpCode != opcode.Query {
		return NoImplMessage
	}

	// If there is more than one question, we reject. Most
	// DNS Servers and resolvers don't implement this
	// feature
	if h.QDCount != 1 {
		return RejectMessage
	}

	// Potentially some more early exist conditions
	return AcceptMessage
}

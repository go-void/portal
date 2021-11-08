package server

import (
	"github.com/go-void/portal/pkg/types/dns"
	"github.com/go-void/portal/pkg/types/opcode"
)

type AcceptFunc func(dns.Header) AcceptAction

type AcceptAction int

const (
	AcceptMessage AcceptAction = iota
	RejectMessage
	IgnoreMessage
	NoImplMessage
)

func DefaultAcceptFunc(h dns.Header) AcceptAction {
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

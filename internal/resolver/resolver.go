package resolver

import (
	"errors"

	"github.com/go-void/portal/internal/types/dns"
	"github.com/go-void/portal/internal/types/rr"
)

var (
	ErrNoAnswer = errors.New("no answer")
)

// Resolver describes a DNS resolver which can resolve
// DNS queries (of unknown / external names) iteratively
// or recursively
type Resolver interface {
	// Resolve resolves a query by first looking if
	// there is a cached record and if not it will
	// look up using remote DNS servers
	Resolve(string, uint16, uint16) (rr.RR, error)

	// ResolveQuestion is a convenience function which
	// allows to provide a DNS question instead of
	// individual parameters
	ResolveQuestion(dns.Question) (rr.RR, error)

	// Lookup looks up a query either iteratively,
	// recursively or by forwarding it to a upstream
	// DNS server
	Lookup(string, uint16, uint16) (rr.RR, error)

	// Refresh refreshes a cached record by looking
	// it up again
	Refresh(string, uint16, uint16)
}

// TODO (Techassi): Figure this out
type ResolveChain struct {
	QueryName string
	// Links []Link
}

package resolver

import (
	"github.com/go-void/portal/internal/cache"
	"github.com/go-void/portal/internal/types/dns"
	"github.com/go-void/portal/internal/types/rr"
)

// Resolver describes a DNS resolver which can resolve
// DNS queries (of unknown / external names) iteratively
// or recursively
type Resolver interface {
	// Resolve resolves a query either iteratively,
	// recursively or by forwarding it to a upstream
	// DNS server
	Resolve(string, uint16, uint16) (rr.RR, error)

	// ResolveQuestion is a convenience function which
	// allows to provide a DNS question instead of
	// individual parameters
	ResolveQuestion(dns.Question) (rr.RR, error)

	// Cache provides access to a cache instance to
	// store answers from remote DNS servers for
	// the given TTL
	Cache() cache.Cache
}

// TODO (Techassi): Figure this out
type ResolveChain struct {
	QueryName string
	// Links []Link
}

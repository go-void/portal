package resolver

import (
	"errors"

	"github.com/go-void/portal/pkg/cache"
	"github.com/go-void/portal/pkg/config"
	"github.com/go-void/portal/pkg/types/dns"
	"github.com/go-void/portal/pkg/types/rr"
)

var (
	ErrFatal    = errors.New("resolver: fatal error")
	ErrNoAnswer = errors.New("resolver: no answer")
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

func New(cfg config.ResolverOptions, c cache.Cache) Resolver {
	switch cfg.Mode {
	case "r":
		return NewRecursiveResolver(cfg, c)
	case "i":
		return NewIterativeResolver()
	case "f":
		return NewForwardingResolver(cfg.Upstream, c)
	}
	return nil
}

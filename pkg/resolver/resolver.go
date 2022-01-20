package resolver

import (
	"errors"

	"github.com/go-void/portal/pkg/cache"
	"github.com/go-void/portal/pkg/config"
	"github.com/go-void/portal/pkg/logger"
	"github.com/go-void/portal/pkg/types/dns"
)

var (
	ErrFatal    = errors.New("resolver: fatal error")
	ErrNoAnswer = errors.New("resolver: no answer")
)

// TODO (Techassi): Rework resolver API to more easily add one or more RRs to any RR section

// Resolver describes a DNS resolver which can resolve DNS queries (of
// unknown / external names) iteratively or recursively
type Resolver interface {
	// Resolve resolves a query of a DNS message by first looking up the
	// request in the cache and if not found will continue to lookup via
	// the Lookup function
	Resolve(*dns.Message) (Result, error)

	ResolveRaw(string, uint16, uint16) (Result, error)

	// Lookup looks up a query either iteratively, recursively or by
	// forwarding it to a upstream DNS server
	Lookup(string, uint16, uint16) (Result, error)

	// Refresh refreshes a cached record by looking it up again
	Refresh(string, uint16, uint16)
}

// TODO (Techassi): Figure this out
type ResolveChain struct {
	QueryName string
	// Links []Link
}

func New(cfg config.ResolverOptions, c cache.Cache, l *logger.Logger) Resolver {
	switch cfg.Mode {
	case "r":
		return NewRecursiveResolver(cfg, c, l)
	case "i":
		// return NewIterativeResolver()
	case "f":
		return NewForwardingResolver(cfg, c, l)
	}
	return nil
}

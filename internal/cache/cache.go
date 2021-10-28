// Package cache provides an in-memory cache to store
// answers from remote DNS servers
package cache

import (
	"github.com/go-void/portal/internal/types/dns"
	"github.com/go-void/portal/internal/types/rr"
)

type Cache interface {
	Get(dns.Question) (rr.RR, bool)

	Set(string, string)

	Len() int
}

package resolver

import (
	"errors"
	"net"
	"sync"

	"github.com/go-void/portal/internal/cache"
	"github.com/go-void/portal/internal/client"
	"github.com/go-void/portal/internal/types/dns"
	"github.com/go-void/portal/internal/types/rr"
)

var (
	ErrNoData = errors.New("no data")
)

type RecursiveResolver struct {
	// Client is a DNS client which sends queries to
	// external DNS servers
	Client client.Client

	// hints is a slice of root DNS server hints
	Hints []net.IP

	// hintIndex keeps track of which root server should
	// be used. It is a simple round-robin algorithm
	HintIndex int

	// Access to the cache instance
	cache cache.Cache

	lock sync.RWMutex
}

// NewRecursiveResolver returns a new recursive resolver
func NewRecursiveResolver(hints []net.IP) *RecursiveResolver {
	return &RecursiveResolver{
		Client: client.NewDefaultClient(),
		Hints:  hints,
	}
}

// Resolve resolves a query by recursivly resolving it
func (r *RecursiveResolver) Resolve(name string, class, t uint16) (rr.RR, error) {
	var glueFound bool
	var ip = r.Hint()

	for {
		response, err := r.Client.Query(name, class, t, ip)
		if err != nil {
			return nil, err
		}

		// We got a answer, return it immediatly
		if response.Header.ANCount > 0 {
			return response.Answer[0], nil
		}

		// We got referals to remote DNS servers
		if response.Header.NSCount > 0 {
			// NOTE (Techassi): Should we check all NS records?
			n := response.Authority[0]

			// NOTE (Techassi): Can we be sure that each record in the authority section is a NS record?
			ns := n.(*rr.NS)

			// Check if we have a "glue record" for the NS record
			if response.Header.ARCount > 0 {
				glueFound = false

				for _, ar := range response.Additional {
					if ar.Header().Class != class || ar.Header().Type != t {
						continue
					}

					if ar.Header().Name == ns.NSDName {
						// At this point we found a glue record

						switch ar.Header().Type {
						case rr.TypeA:
							glue := ar.(*rr.A)
							ip = glue.Address
						case rr.TypeAAAA:
							glue := ar.(*rr.AAAA)
							ip = glue.Address
						}

						glueFound = true
						break
					}
				}

				// We have a glue record, which means we can start from the top and use the selected IP address to
				// continue our recursive query
				if glueFound {
					continue
				}
			}

			// We don't have any matching "glue record", so we have to manually lookup the NS record domain name
			nsAnswer, err := r.Resolve(ns.NSDName, 1, 1)
			if err != nil {
				return nil, err
			}

			// At this point we should have an A record for the remote authoritative name server
			nsARecord := nsAnswer.(*rr.A)
			ip = nsARecord.Address
			continue
		}

		// If we are here we did get a response but with no content
		return nil, ErrNoData
	}
}

// ResolveQuestion is a convenience function which allows to provide a DNS question instead of individual parameters
func (r *RecursiveResolver) ResolveQuestion(question dns.Question) (rr.RR, error) {
	return r.Resolve(question.Name, question.Class, question.Type)
}

// Cache provides access to a cache instance to store answers from remote DNS servers for the given TTL
func (r *RecursiveResolver) Cache() cache.Cache {
	return r.cache
}

// Hint returns a root hint
func (r *RecursiveResolver) Hint() net.IP {
	r.lock.Lock()

	if r.HintIndex == len(r.Hints)-1 {
		r.HintIndex = 0
	} else {
		r.HintIndex++
	}
	i := r.HintIndex

	r.lock.Unlock()
	return r.Hints[i]
}

package resolver

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/go-void/portal/pkg/cache"
	"github.com/go-void/portal/pkg/client"
	"github.com/go-void/portal/pkg/types/dns"
	"github.com/go-void/portal/pkg/types/rr"
)

type RecursiveResolver struct {
	// Client is a DNS client which sends queries to
	// external DNS servers
	Client client.Client

	MaxExpired int

	// hints is a slice of root DNS server hints
	Hints []net.IP

	// HintIndex keeps track of which root server should
	// be used. It is a simple round-robin algorithm
	HintIndex int

	// Access to the cache instance
	Cache cache.Cache

	lock sync.RWMutex
}

// NewRecursiveResolver returns a new recursive resolver
func NewRecursiveResolver(hints []net.IP, c cache.Cache) *RecursiveResolver {
	return &RecursiveResolver{
		Client: client.NewDefault(),
		Hints:  hints,
		Cache:  c,
	}
}

// Resolve resolves a query by recursivly resolving it
func (r *RecursiveResolver) Resolve(name string, class, t uint16) (rr.RR, error) {
	entry, status, err := r.Cache.Lookup(name, class, t)
	if err != nil {
		fmt.Println(err)
	}

	if status == cache.Hit {
		return entry.Record, nil
	}

	if status == cache.Expired {
		max := entry.Expire.Add(time.Duration(r.MaxExpired) * time.Second)
		if max.After(time.Now()) {
			go r.Refresh(name, class, t)
			return entry.Record, nil
		}
	}

	response, err := r.Lookup(name, class, t)
	if err != nil {
		return nil, err
	}

	err = r.Cache.Set(name, class, t, response, response.Header().TTL)
	if err != nil {
		fmt.Println(err)
	}

	return response, nil
}

// ResolveQuestion is a convenience function which allows to provide a DNS question instead of individual parameters
func (r *RecursiveResolver) ResolveQuestion(question dns.Question) (rr.RR, error) {
	return r.Resolve(question.Name, question.Class, question.Type)
}

func (r *RecursiveResolver) Lookup(name string, class, t uint16) (rr.RR, error) {
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
		return nil, ErrNoAnswer
	}
}

func (r *RecursiveResolver) Refresh(name string, class, t uint16) {
	return
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

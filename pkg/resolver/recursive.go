package resolver

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/go-void/portal/pkg/cache"
	"github.com/go-void/portal/pkg/client"
	"github.com/go-void/portal/pkg/config"
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

	cacheEnabled bool
	lock         sync.RWMutex
}

// NewRecursiveResolver returns a new recursive resolver
func NewRecursiveResolver(cfg config.ResolverOptions, c cache.Cache) *RecursiveResolver {
	// TODO (Techassi): Read in hints
	hints := []net.IP{
		net.ParseIP("198.41.0.4"),
		net.ParseIP("199.9.14.201"),
	}

	return &RecursiveResolver{
		Client:       client.NewDefault(),
		Hints:        hints,
		Cache:        c,
		cacheEnabled: cfg.CacheEnabled,
	}
}

// ResolveQuestion is a convenience function which allows to provide a DNS question instead of individual parameters
func (r *RecursiveResolver) ResolveQuestion(question dns.Question) (rr.RR, error) {
	return r.Resolve(question.Name, question.Class, question.Type)
}

// Resolve resolves a query by recursivly resolving it
func (r *RecursiveResolver) Resolve(name string, class, t uint16) (rr.RR, error) {
	if r.cacheEnabled {
		record, ok := r.LookupInCache(name, class, t)
		if ok {
			return record, nil
		}
	}

	response, err := r.Lookup(name, class, t)
	if err != nil {
		return nil, err
	}

	if r.cacheEnabled {
		err = r.Cache.Set(name, class, t, response, response.Header().TTL)
		if err != nil {
			fmt.Println(err)
		}
	}

	return response, nil
}

func (r *RecursiveResolver) Lookup(name string, class, t uint16) (rr.RR, error) {
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

		// We have no direct answer or references to
		// authoriative DNS servers
		if response.Header.NSCount == 0 {
			return nil, ErrNoAnswer
		}

		// Try to find glue records
		newIP, err := r.FindGlue(response, class, t)
		if err != nil {
			return nil, err
		}
		ip = newIP
	}
}

func (r *RecursiveResolver) Refresh(name string, class, t uint16) {
	response, err := r.Lookup(name, class, t)
	if err != nil {
		// NOTE (Techassi): Log this
		return
	}

	err = r.Cache.Set(name, class, t, response, response.Header().TTL)
	if err != nil {
		// NOTE (Techassi): Log this
	}
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

// LookupInCache is a convenience function which abstracts the lookup of a domain name in the cache
func (r *RecursiveResolver) LookupInCache(name string, class, t uint16) (rr.RR, bool) {
	entry, status, err := r.Cache.Lookup(name, class, t)
	if err != nil {
		fmt.Println(err)
		return nil, false
	}

	if status == cache.Hit {
		return entry.Record, true
	}

	if status == cache.Expired {
		max := entry.Expire.Add(time.Duration(r.MaxExpired) * time.Second)
		if max.After(time.Now()) {
			go r.Refresh(name, class, t)
			return entry.Record, true
		}
	}

	return nil, false
}

// FindGlue tries to find glue records for provided authoriative DNS servers
func (r *RecursiveResolver) FindGlue(response dns.Message, class, t uint16) (net.IP, error) {
	for _, nsrr := range response.Authority {
		// NOTE (Techassi): Can we be sure that each record in the authority section is a NS record?
		ns := nsrr.(*rr.NS)

		for _, adrr := range response.Additional {
			if adrr.Header().Class != class || adrr.Header().Type != t {
				continue
			}

			if adrr.Header().Name != ns.NSDName {
				continue
			}

			var ip net.IP

			switch adrr.Header().Type {
			case rr.TypeA:
				ip = adrr.(*rr.A).Address
			case rr.TypeAAAA:
				ip = adrr.(*rr.AAAA).Address
			}

			// Cache NS records
			if r.cacheEnabled {
				err := r.Cache.Set(ns.NSDName, class, t, adrr, adrr.Header().TTL)
				if err != nil {
					// TODO (Techassi): Handle error
					fmt.Println(err)
				}
			}

			return ip, nil
		}

		// We don't have any matching "glue record", so we have to manually lookup the NS record domain name
		nsAnswer, err := r.Resolve(ns.NSDName, 1, 1)
		if err != nil {
			// Log error
			continue
		}

		// At this point we should have an A record for the remote authoritative name server
		nsARecord := nsAnswer.(*rr.A)
		ip := nsARecord.Address
		return ip, nil
	}

	return nil, ErrFatal
}

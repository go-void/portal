package resolver

import (
	"fmt"
	"net/netip"
	"sync"

	"github.com/go-void/portal/pkg/cache"
	"github.com/go-void/portal/pkg/client"
	"github.com/go-void/portal/pkg/config"
	"github.com/go-void/portal/pkg/logger"
	"github.com/go-void/portal/pkg/types/dns"
	"github.com/go-void/portal/pkg/types/rr"
)

type RecursiveResolver struct {
	// client is a DNS client which sends queries to
	// external DNS servers
	client *client.Client

	// Access to the logger instance
	logger *logger.Logger

	// Access to the cache instance
	cache cache.Cache

	maxExpired int

	// hints is a slice of root DNS server hints
	hints []netip.Addr

	// hintIndex keeps track of which root server should
	// be used. It is a simple round-robin algorithm
	hintIndex int

	cacheEnabled bool
	lock         sync.RWMutex
}

// NewRecursiveResolver returns a new recursive resolver
func NewRecursiveResolver(cfg config.ResolverOptions, c cache.Cache, l *logger.Logger) *RecursiveResolver {
	// TODO (Techassi): Read in hints
	hints := []netip.Addr{
		netip.AddrFrom4([4]byte{198, 41, 0, 4}),
		netip.AddrFrom4([4]byte{199, 9, 14, 201}),
	}

	return &RecursiveResolver{
		client:       client.New(l),
		hints:        hints,
		cache:        c,
		logger:       l,
		cacheEnabled: cfg.CacheEnabled,
	}
}

// Resolve recursivly resolves a query
func (r *RecursiveResolver) Resolve(message *dns.Message) (Result, error) {
	name, class, t := message.Q()
	return r.ResolveRaw(name, class, t)
}

// TODO (Techassi): We have duplicate code. That's meh... Fix it future me!
func (r *RecursiveResolver) ResolveRaw(name string, class, t uint16) (Result, error) {
	if r.cacheEnabled {
		records, ok := r.LookupInCache(name, class, t)
		if ok {
			// NOTE (Techassi): This is not the correct way to do it
			return Result{
				Answer: records,
			}, nil
		}
	}

	result, err := r.Lookup(name, class, t)
	if err != nil {
		return result, err
	}

	if r.cacheEnabled {
		// TODO (Techassi): Figure out a simple, elegant and performant way of setting ALL records of result
		err = r.cache.Set(name, result.Answer)
		if err != nil {
			fmt.Println(err)
		}
	}

	return result, nil
}

func (r *RecursiveResolver) Lookup(name string, class, t uint16) (Result, error) {
	var ip = r.Hint()

	for {
		response, err := r.client.Query(name, class, t, ip)
		if err != nil {
			return Result{}, err
		}

		// We got an answer, return it immediatly
		if response.Header.ANCount > 0 {
			return NewResult(response), nil
		}

		// We have no direct answer or references to
		// authoriative DNS servers
		if response.Header.NSCount == 0 {
			return Result{}, ErrNoAnswer
		}

		// TODO (Techassi): Figure out a way how we should handle SOA records
		//                  at this stage as there is no way possible to find
		//                  any glue records. We instead should just return
		//                  the received response. Maybe introduce the "error"
		//                  NoErrSOAFound to indicate the glue search stopped
		//                  because we found a SOA record.

		if response.IsSOA() {
			return NewResult(response), nil
		}

		// Try to find glue records
		newIP, err := r.FindGlue(response, class, t)
		if err != nil {
			return Result{}, err
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

	err = r.cache.Set(name, response.Answer)
	if err != nil {
		// NOTE (Techassi): Log this
	}
}

// Hint returns a root hint
func (r *RecursiveResolver) Hint() netip.Addr {
	r.lock.Lock()

	if r.hintIndex == len(r.hints)-1 {
		r.hintIndex = 0
	} else {
		r.hintIndex++
	}
	i := r.hintIndex

	r.lock.Unlock()
	return r.hints[i]
}

// LookupInCache is a convenience function which abstracts the lookup of a domain name in the cache
func (r *RecursiveResolver) LookupInCache(name string, class, t uint16) ([]rr.RR, bool) {
	records, status, err := r.cache.Lookup(name, class, t)
	if err != nil {
		fmt.Println(err)
		return nil, false
	}

	if status == cache.Hit {
		return records, true
	}

	// TODO (Techassi): Redo this
	// if status == cache.Expired {
	// 	max := entry.Expire.Add(time.Duration(r.maxExpired) * time.Second)
	// 	if max.After(time.Now()) {
	// 		go r.Refresh(name, class, t)
	// 		return entry.Record, true
	// 	}
	// }

	return nil, false
}

// NOTE (Techassi): Should we be able to return multiple IPs here?
// FindGlue tries to find glue records for provided authoriative DNS servers
func (r *RecursiveResolver) FindGlue(response *dns.Message, class, t uint16) (netip.Addr, error) {
	for _, nsrr := range response.Authority {
		ns, ok := nsrr.(*rr.NS)
		if !ok {
			continue
		}

		for _, adrr := range response.Additional {
			if adrr.Header().Class != class || adrr.Header().Type != t {
				continue
			}

			if adrr.Header().Name != ns.NSDName {
				continue
			}

			var ip netip.Addr

			switch adrr.Header().Type {
			case rr.TypeA:
				ip = adrr.(*rr.A).Address
			case rr.TypeAAAA:
				ip = adrr.(*rr.AAAA).Address
			}

			// Cache NS records
			// if r.cacheEnabled {
			// 	err := r.cache.Set(ns.NSDName, class, t, adrr)
			// 	if err != nil {
			// 		// TODO (Techassi): Handle error
			// 		fmt.Println(err)
			// 	}
			// }

			return ip, nil
		}

		// We don't have any matching "glue record", so we have to manually lookup the NS record domain name
		result, err := r.ResolveRaw(ns.NSDName, 1, 1)
		if err != nil {
			// Log error
			continue
		}

		// At this point we should have an A record for the remote authoritative name server
		nsARecord := result.Answer[0].(*rr.A)
		ip := nsARecord.Address
		return ip, nil
	}

	return netip.Addr{}, ErrFatal
}

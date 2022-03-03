package resolver

import (
	"fmt"
	"net"

	"github.com/go-void/portal/pkg/cache"
	"github.com/go-void/portal/pkg/client"
	"github.com/go-void/portal/pkg/config"
	"github.com/go-void/portal/pkg/logger"
	"github.com/go-void/portal/pkg/types/dns"
	"github.com/go-void/portal/pkg/types/rr"
)

type ForwardingResolver struct {
	// client is a DNS client which sends queries to
	// external DNS servers
	client *client.Client

	// Access to the cache instance
	cache cache.Cache

	// Access to the logger instance
	logger *logger.Logger

	// upstream is a IP address of a upstream DNS
	// server
	upstream net.IP

	maxExpired int

	cacheEnabled bool
}

// NewForwardingResolver returns a new forwarding resolver
func NewForwardingResolver(cfg config.ResolverOptions, c cache.Cache, l *logger.Logger) *ForwardingResolver {
	return &ForwardingResolver{
		client:       client.New(l),
		upstream:     net.ParseIP(cfg.RawUpstream),
		cacheEnabled: cfg.CacheEnabled,
		maxExpired:   300,
		cache:        c,
		logger:       l,
	}
}

// Resolve resolves a query by forwarding it to the upstream DNS server
func (r *ForwardingResolver) Resolve(message *dns.Message) (Result, error) {
	name, class, t := message.Q()
	return r.ResolveRaw(name, class, t)
}

func (r *ForwardingResolver) ResolveRaw(name string, class, t uint16) (Result, error) {
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

func (r *ForwardingResolver) Lookup(name string, class, t uint16) (Result, error) {
	response, err := r.client.Query(name, class, t, r.upstream)
	if err != nil {
		return Result{}, err
	}

	return NewResult(response), nil
}

func (r *ForwardingResolver) Refresh(name string, class, t uint16) {
	response, err := r.Lookup(name, class, t)
	if err != nil {
		// NOTE (Techassi): Log this
		return
	}

	// TODO (Techassi): Not just the answer
	err = r.cache.Set(name, response.Answer)
	if err != nil {
		// NOTE (Techassi): Log this
	}
}

// LookupInCache is a convenience function which abstracts the lookup of a domain name in the cache
func (r *ForwardingResolver) LookupInCache(name string, class, t uint16) ([]rr.RR, bool) {
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

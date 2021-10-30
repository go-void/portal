package resolver

import (
	"fmt"
	"net"
	"time"

	"github.com/go-void/portal/internal/cache"
	"github.com/go-void/portal/internal/client"
	"github.com/go-void/portal/internal/types/dns"
	"github.com/go-void/portal/internal/types/rr"
)

type ForwardingResolver struct {
	// Client is a DNS client which sends queries to
	// external DNS servers
	Client client.Client

	// Upstream is a IP address of a upstream DNS
	// server
	Upstream net.IP

	MaxExpired int

	Cache cache.Cache
}

// NewForwardingResolver returns a new forwarding resolver
func NewForwardingResolver(upstream net.IP, c cache.Cache) *ForwardingResolver {
	return &ForwardingResolver{
		Client:   client.NewDefaultClient(),
		Upstream: upstream,
		Cache:    c,
	}
}

// Resolve resolves a query by forwarding it to the upstream DNS server
func (r *ForwardingResolver) Resolve(name string, class, t uint16) (rr.RR, error) {
	entry, status, err := r.Cache.Get(name, class, t)
	if err != nil {
		return nil, err
	}

	fmt.Println("Cache status:", status)

	if status == cache.Hit {
		return entry.Data, nil
	}

	if status == cache.Expired {
		max := time.Now().Add(time.Duration(r.MaxExpired) * time.Second)
		if max.Sub(entry.Expire) >= 0 {
			// Return expired record but fetch again
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

func (r *ForwardingResolver) ResolveQuestion(question dns.Question) (rr.RR, error) {
	return r.Resolve(question.Name, question.Class, question.Type)
}

func (r *ForwardingResolver) Lookup(name string, class, t uint16) (rr.RR, error) {
	response, err := r.Client.Query(name, class, t, r.Upstream)
	if err != nil {
		return nil, err
	}

	if len(response.Answer) <= 0 {
		return nil, ErrNoAnswer
	}

	return response.Answer[0], nil
}

func (r *ForwardingResolver) Refresh(name string, class, t uint16) {
	return
}

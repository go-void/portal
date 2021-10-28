package resolver

import (
	"net"

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

	cache cache.Cache
}

// NewForwardingResolver returns a new forwarding resolver
func NewForwardingResolver(upstream net.IP) *ForwardingResolver {
	return &ForwardingResolver{
		Client:   client.NewDefaultClient(),
		Upstream: upstream,
	}
}

// Resolve resolves a query by forwarding it to the upstream DNS server
func (r *ForwardingResolver) Resolve(question dns.Question) (rr.RR, error) {
	response, err := r.Client.Query(question.Name, question.Class, question.Type, r.Upstream)
	if err != nil {
		return nil, err
	}
	return response.Answer[0], nil
}

func (r *ForwardingResolver) Cache() cache.Cache {
	return r.cache
}

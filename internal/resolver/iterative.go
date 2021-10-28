package resolver

import (
	"net"

	"github.com/go-void/portal/internal/cache"
	"github.com/go-void/portal/internal/client"
	"github.com/go-void/portal/internal/types/dns"
	"github.com/go-void/portal/internal/types/rr"
)

type IterativeResolver struct {
	// Client is a DNS client which sends queries to
	// external DNS servers
	Client client.Client

	cache cache.Cache
}

// NewIterativeResolver returns a new iterative resolver
func NewIterativeResolver() *IterativeResolver {
	return &IterativeResolver{
		Client: client.NewDefaultClient(),
	}
}

// Resolve resolves a query by answering with a DNS server which is closer to the quried name
func (r *IterativeResolver) Resolve(question dns.Question) (rr.RR, error) {
	// TODO (Techassi): Which IP should we use here?
	response, err := r.Client.Query(question.Name, question.Class, question.Type, net.IP{})
	if err != nil {
		return nil, err
	}
	return response.Answer[0], nil
}

func (r *IterativeResolver) Cache() cache.Cache {
	return r.cache
}

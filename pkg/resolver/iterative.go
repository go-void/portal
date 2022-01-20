package resolver

import (
	"github.com/go-void/portal/pkg/cache"
	"github.com/go-void/portal/pkg/client"
	"github.com/go-void/portal/pkg/types/dns"
	"github.com/go-void/portal/pkg/types/rr"
)

type IterativeResolver struct {
	// Client is a DNS client which sends queries to
	// external DNS servers
	Client client.Client

	Cache cache.Cache
}

// NewIterativeResolver returns a new iterative resolver
func NewIterativeResolver() *IterativeResolver {
	return &IterativeResolver{
		// Client: client.NewDefault(),
	}
}

// Resolve resolves a query by answering with a DNS server which is closer to the quried name
func (r *IterativeResolver) Resolve(name string, class, t uint16) (rr.RR, error) {
	// TODO (Techassi): Create answer with referals to external DNS servers
	return nil, nil
}

func (r *IterativeResolver) ResolveQuestion(question dns.Question) (rr.RR, error) {
	return r.Resolve(question.Name, question.Class, question.Type)
}

func (r *IterativeResolver) Lookup(name string, class, t uint16) (rr.RR, error) {
	return nil, nil
}

func (r *IterativeResolver) Refresh(name string, class, t uint16) {}

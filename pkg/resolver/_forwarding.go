package resolver

import (
	"fmt"
	"net"
	"time"

	"github.com/go-void/portal/pkg/cache"
	"github.com/go-void/portal/pkg/client"
	"github.com/go-void/portal/pkg/logger"
	"github.com/go-void/portal/pkg/types/dns"
	"github.com/go-void/portal/pkg/types/rr"
)

type ForwardingResolver struct {
	// client is a DNS client which sends queries to
	// external DNS servers
	client client.Client

	// upstream is a IP address of a upstream DNS
	// server
	upstream net.IP

	maxExpired int

	cache  cache.Cache
	logger *logger.Logger
}

// NewForwardingResolver returns a new forwarding resolver
func NewForwardingResolver(upstream net.IP, c cache.Cache, l *logger.Logger) *ForwardingResolver {
	return &ForwardingResolver{
		client:     client.NewDefault(),
		upstream:   upstream,
		maxExpired: 300,
		cache:      c,
		logger:     l,
	}
}

// Resolve resolves a query by forwarding it to the upstream DNS server
func (r *ForwardingResolver) Resolve(name string, class, t uint16) ([]rr.RR, error) {
	entry, status, err := r.cache.Lookup(name, class, t)
	if err != nil {
		fmt.Println(err)
	}

	if status == cache.Hit {
		return entry, nil
	}

	// TODO (Techassi): Redo this
	// if status == cache.Expired {
	// 	max := entry.Expire.Add(time.Duration(r.maxExpired) * time.Second)
	// 	if max.After(time.Now()) {
	// 		go r.Refresh(name, class, t)
	// 		return entry.Record, true
	// 	}
	// }

	response, err := r.Lookup(name, class, t)
	if err != nil {
		return nil, err
	}

	err = r.cache.Set(name, class, t, response, response.Header().TTL)
	if err != nil {
		fmt.Println(err)
	}

	return response, nil
}

func (r *ForwardingResolver) ResolveQuestion(question dns.Question) (rr.RR, error) {
	return r.Resolve(question.Name, question.Class, question.Type)
}

func (r *ForwardingResolver) Lookup(name string, class, t uint16) ([]rr.RR, error) {
	response, err := r.client.Query(name, class, t, r.upstream)
	if err != nil {
		return nil, err
	}

	if len(response.Answer) <= 0 {
		return nil, ErrNoAnswer
	}

	return response.Answer[0], nil
}

func (r *ForwardingResolver) Refresh(name string, class, t uint16) {
	response, err := r.Lookup(name, class, t)
	if err != nil {
		// NOTE (Techassi): Log this
		return
	}

	err = r.cache.Set(name, class, t, response, response.Header().TTL)
	if err != nil {
		// NOTE (Techassi): Log this
	}
}

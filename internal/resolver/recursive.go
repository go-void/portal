package resolver

import (
	"fmt"
	"net"
	"sync"

	"github.com/go-void/portal/internal/cache"
	"github.com/go-void/portal/internal/client"
	"github.com/go-void/portal/internal/labels"
	"github.com/go-void/portal/internal/types/dns"
	"github.com/go-void/portal/internal/types/rr"
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
func (r *RecursiveResolver) Resolve(question dns.Question) (rr.RR, error) {
	lbs := labels.FromRoot(question.Name)[1:]

	for i, l := range lbs {
		if i == 0 {
			response, err := r.Client.Query(l, question.Class, question.Type, r.Hint())
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			fmt.Println("Resp", response)
			continue
		}

		// Do iterative queries based on previous answers
		// response, err := r.client.Query(l, class, t, r.Hint())
		// if err != nil {
		// 	return nil, err
		// }
		// fmt.Println(response)
	}

	return nil, nil
}

func (r *RecursiveResolver) Cache() cache.Cache {
	return r.cache
}

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

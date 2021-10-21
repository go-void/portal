package resolver

import (
	"sync"

	"github.com/go-void/portal/internal/client"
)

// Resolver describes a DNS resolver which can resolve
// DNS queries (of unknown / external names) iteratively
// or recursively
type Resolver interface {
	// Resolve resolves a query either iteratively
	// or recursively, which is defined by the RD
	// header field
	Resolve(bool, string, uint16, uint16)

	// Iterate resolves a name iteratively
	Iterate(string, uint16, uint16)

	// Recurse resolves a name recursively
	Recurse(string, uint16, uint16)

	// Hints returns root hints as a map
	Hints()
}

type DefaultResolver struct {
	Client client.Client

	hints []string

	// hintIndex keeps track of which root server should
	// be used. It is a simple round-robin algorithm
	hintIndex int

	lock sync.RWMutex
}

func NewDefaultResolver() *DefaultResolver {
	return &DefaultResolver{
		Client: client.NewDefaultClient(),
	}
}

func (r *DefaultResolver) Resolve(rd bool, name string, class, t uint16) {
	if rd {
		r.Recurse(name, class, t)
	} else {
		r.Iterate(name, class, t)
	}
}

// Iterate resolves a name iteratively
func (r *DefaultResolver) Iterate(name string, class, t uint16) {
	// names := utils.LabelsFromRoot(name)[1:]
	// index := r.getHintIndex()
	// fmt.Println(names, index)
}

// Recurse resolves a name recursively
func (r *DefaultResolver) Recurse(name string, class, t uint16) {
	panic("not implemented") // TODO: Implement
}

// Hints returns root hints as a map
func (r *DefaultResolver) Hints() {
	panic("not implemented") // TODO: Implement
}

func (r *DefaultResolver) getHintIndex() int {
	r.lock.RLock()
	i := r.hintIndex
	r.lock.RUnlock()
	return i
}

package resolver

import (
	"errors"
	"net"
	"sync"

	"github.com/go-void/portal/internal/client"
	"github.com/go-void/portal/internal/types/rr"
)

type Mode int

const (
	Recurse Mode = iota
	Iterate
	Forward
)

var (
	ErrInvalidMode = errors.New("invalid mode")
)

// Resolver describes a DNS resolver which can resolve
// DNS queries (of unknown / external names) iteratively
// or recursively
type Resolver interface {
	// Resolve resolves a query either iteratively,
	// recursively or by forwarding it to a upstream
	// DNS server
	Resolve(Mode, string, uint16, uint16) (rr.RR, error)

	// Iterate resolves a name iteratively
	Iterate(string, uint16, uint16) (rr.RR, error)

	// Recurse resolves a name recursively
	Recurse(string, uint16, uint16) (rr.RR, error)

	// Forward forwards a query to a upstream DNS server
	// like 1.1.1.1, 8.8.8.8 or 9.9.9.9
	Forward(string, uint16, uint16) (rr.RR, error)

	// Hints returns root hints as a map
	Hints()
}

type DefaultResolver struct {
	// client is a DNS client which sends queries to
	// external DNS servers
	client client.Client

	// upstream is a IP address of a upstream DNS
	// server
	upstream net.IP

	// hints is a slice of root DNS server hints
	hints []string

	// hintIndex keeps track of which root server should
	// be used. It is a simple round-robin algorithm
	hintIndex int

	lock sync.RWMutex
}

// TODO (Techassi): Figure this out
type ResolveChain struct {
	QueryName string
	// Links []Link
}

func NewDefaultResolver() *DefaultResolver {
	return &DefaultResolver{
		client:   client.NewDefaultClient(),
		upstream: net.ParseIP("1.1.1.1"), // TODO (Techassi): Make this configurable
	}
}

// Resolve resolves a query either iteratively, recursively or by forwarding it to a
// upstream DNS server
func (r *DefaultResolver) Resolve(mode Mode, name string, class, t uint16) (rr.RR, error) {
	switch mode {
	case Recurse:
		return r.Recurse(name, class, t)
	case Iterate:
		return r.Iterate(name, class, t)
	case Forward:
		return r.Forward(name, class, t)
	default:
		return nil, ErrInvalidMode
	}
}

// Iterate resolves a name iteratively
func (r *DefaultResolver) Iterate(name string, class, t uint16) (rr.RR, error) {
	// names := utils.LabelsFromRoot(name)[1:]
	// index := r.getHintIndex()
	// fmt.Println(names, index)
	return nil, nil
}

// Recurse resolves a name recursively
func (r *DefaultResolver) Recurse(name string, class, t uint16) (rr.RR, error) {
	panic("not implemented") // TODO: Implement
}

func (r *DefaultResolver) Forward(name string, class, t uint16) (rr.RR, error) {
	response, err := r.client.Query(name, class, t, r.upstream)
	if err != nil {
		return nil, err
	}
	return response.Answer[0], nil
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

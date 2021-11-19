// Package store provides a domain space data store to retrieve
// and store such data
package store

import (
	"time"

	"github.com/go-void/portal/pkg/tree"
	"github.com/go-void/portal/pkg/types/dns"
	"github.com/go-void/portal/pkg/types/rr"
)

type Store interface {
	// Get matches a 'node' by name and returns the selected type's data.
	// Example: example.com with type A would return 93.184.216.34
	Get(dns.Question) (rr.RR, error)

	// Set sets type's data of a 'node' selected by name
	Set(string, uint16, uint16, rr.RR, uint32) error

	// Indicates if this store is using a cache. This is especially
	// usefull when the store itself is in-memory which eliminates
	// the need of a in-memory cache
	UsesCache() bool
}

// DefaultStore implements a default store based on a in-memory tree structure
type DefaultStore struct {
	*tree.Tree
}

func NewDefault() *DefaultStore {
	return &DefaultStore{
		Tree: tree.New(),
	}
}

// Get matches a 'node' by name and returns the selected type's data
func (s *DefaultStore) Get(question dns.Question) (rr.RR, error) {
	node, err := s.Tree.Get(question.Name)
	if err != nil {
		return nil, err
	}

	nodeRecord, err := node.Entry(question.Class, question.Type)
	if err != nil {
		return nil, err
	}

	return nodeRecord.Record, err
}

// Set sets type's data of a 'node' selected by name
func (s *DefaultStore) Set(name string, class, t uint16, record rr.RR, ttl uint32) error {
	node, err := s.Tree.Populate(name)
	if err != nil {
		return tree.ErrNodeNotFound
	}

	expire := time.Now().Add(time.Duration(ttl) * time.Second)
	node.SetData(class, t, record, expire)

	return nil
}

func (s *DefaultStore) UsesCache() bool {
	return false
}

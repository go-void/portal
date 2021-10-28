// Package store provides a domain space data store to retrieve
// and store such data
package store

import (
	"github.com/go-void/portal/internal/tree"
	"github.com/go-void/portal/internal/types/dns"
	"github.com/go-void/portal/internal/types/rr"
)

type Store interface {
	// Get matches a 'node' by name and returns the selected type's data.
	// Example: example.com with type A would return 93.184.216.34
	Get(dns.Question) (rr.RR, error)

	// Set sets type's data of a 'node' selected by name
	Set(string, uint16, uint16, interface{}) error

	// Indicates if this store is using a cache. This is especially
	// usefull when the store itself is in-memory which eliminates
	// the need of a in-memory cache
	UsesCache() bool
}

// DefaultStore implements a default store based on a in-memory tree
// structure
type DefaultStore struct {
	Tree *tree.Tree
}

func NewDefaultStore() *DefaultStore {
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

	data, err := node.Data(question.Class, question.Type)
	if err != nil {
		return nil, err
	}

	record, err := rr.New(question.Type)
	if err != nil {
		return nil, err
	}

	record.SetHeader(rr.Header{
		Name:     question.Name,
		Class:    question.Class,
		Type:     question.Type,
		TTL:      3600,
		RDLength: 4, // FIXME (Techassi): Don't make this fixed
	})

	err = record.SetData(data)
	if err != nil {
		return nil, err
	}

	return record, err
}

// Set sets type's data of a 'node' selected by name
func (s *DefaultStore) Set(name string, class, t uint16, data interface{}) error {
	node, err := s.Tree.Populate(name)
	if err != nil {
		return tree.ErrNodeNotFound
	}

	node.SetData(class, t, data)
	return nil
}

func (s *DefaultStore) UsesCache() bool {
	return false
}

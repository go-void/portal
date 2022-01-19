// Package store provides an interface to store and retrieve custom DNS records
package store

import (
	"github.com/go-void/portal/pkg/tree"
	"github.com/go-void/portal/pkg/types/dns"
	"github.com/go-void/portal/pkg/types/rr"
)

type Store interface {
	// GetFromQuestion returns records by name and the selected type's data.
	// Example: example.com with type A would return 93.184.216.34
	GetFromQuestion(dns.Question) ([]rr.RR, error)

	// Get returns records by name and the selected type's data.
	// Example: example.com with type A would return 93.184.216.34
	Get(string, uint16, uint16) ([]rr.RR, error)

	// Add adds a new resource records to the store
	Add(string, []rr.RR) error
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

// GetFromQuestion returns a record by name and the selected type's data.
// Example: example.com with type A would return 93.184.216.34
func (s *DefaultStore) GetFromQuestion(question dns.Question) ([]rr.RR, error) {
	node, err := s.Tree.Get(question.Name)
	if err != nil {
		return nil, err
	}

	records, err := node.Records(question.Class, question.Type)
	if err != nil {
		return nil, err
	}

	return records, err
}

// Add adds a new resource record to the store
func (s *DefaultStore) Add(name string, records []rr.RR) error {
	node, err := s.Tree.Populate(name)
	if err != nil {
		return tree.ErrNodeNotFound
	}

	node.AddRecords(records)
	return nil
}

// Get returns a record by name and the selected type's data.
// Example: example.com with type A would return 93.184.216.34
func (s *DefaultStore) Get(name string, class, t uint16) ([]rr.RR, error) {
	node, err := s.Tree.Get(name)
	if err != nil {
		return nil, err
	}

	records, err := node.Records(class, t)
	if err != nil {
		return nil, err
	}

	return records, err
}

// Package cache provides an in-memory cache to store answers from remote DNS servers
package cache

import (
	"errors"
	"time"

	"github.com/go-void/portal/pkg/labels"
	"github.com/go-void/portal/pkg/types/dns"
	"github.com/go-void/portal/pkg/types/rr"
)

var (
	ErrNodeNotFound       = errors.New("cache: node not found in tree")
	ErrNoSuchData         = errors.New("cache: no such data")
	ErrChildAlreadyExists = errors.New("cache: child already exists")
)

// Cache describes a cache interface to store RRs retrieved
// from remote DNS servers
type Cache interface {
	Lookup(string, uint16, uint16) (Entry, Status, error)

	LookupQuestion(dns.Question) (Entry, Status, error)

	Set(string, uint16, uint16, rr.RR, uint32) error

	SetMultiple(string, uint16, uint16, []rr.RR, uint32) error
}

// DefaultCache implements the Cache interface and stores
// RRs in an in-memory tree
type DefaultCache struct {
	root Node
}

// NewDefaultCache returns a new default in-memory tree cache
func NewDefaultCache() *DefaultCache {
	return &DefaultCache{
		root: Node{
			parent:   nil,
			children: make(map[string]Node),
			entries:  make(map[uint16]Entry),
		},
	}
}

// Get walks the tree to retrieve a node (INTERNAL)
func (c *DefaultCache) Get(name string) (Node, error) {
	nodes, err := c.Walk(name)
	if err != nil {
		return Node{}, err
	}
	return nodes[len(nodes)-1], nil
}

// Walk walks the tree until the requested node is found and
// returns errors encountered along the way (INTERNAL)
func (c *DefaultCache) Walk(name string) ([]Node, error) {
	var (
		current = c.root
		nodes   = []Node{}
		names   = labels.FromRoot(name)
	)

	for _, name := range names {
		if name == "" || name == "." {
			nodes = append(nodes, current)
			continue
		}

		node, err := current.Child(name)
		if err != nil {
			return nodes, ErrNodeNotFound
		}

		current = node
		nodes = append(nodes, node)
	}

	return nodes, nil
}

// Populate populates the tree along the provided name (INTERNAL).
// Example: example.com => . -> com -> example
func (c *DefaultCache) Populate(name string) (Node, error) {
	var current = c.root
	var names = labels.FromRoot(name)

	for _, name := range names {
		if name == "" || name == "." {
			continue
		}

		node, err := current.Child(name)
		if err != nil {
			node := Node{
				parent:   &current,
				children: make(map[string]Node),
				entries:  make(map[uint16]Entry),
			}

			err := current.AddChild(name, node)
			if err != nil {
				return Node{}, err
			}

			current = node
			continue
		}
		current = node
	}
	return current, nil
}

// Lookup looks up a entry for name with class and type and returns the status and errors
// encountered along the way
func (c *DefaultCache) Lookup(name string, class, t uint16) (Entry, Status, error) {
	node, err := c.Get(name)
	if err != nil {
		return Entry{}, Miss, nil
	}

	entry, err := node.Entry(class, t)
	if err != nil {
		return Entry{}, Miss, nil
	}

	status := Hit
	if entry.Expire.Before(time.Now()) {
		status = Expired
	}

	rr.UpdateTTL(entry.Record, entry.Expire)

	return Entry{
		Record: entry.Record,
		Expire: entry.Expire,
	}, status, nil
}

// LookupQuestion is a convenience function to lookup a DNS question
func (c *DefaultCache) LookupQuestion(message dns.Question) (Entry, Status, error) {
	return c.Lookup(message.Name, message.Class, message.Type)
}

// Set sets (or adds) a new cache entry
func (c *DefaultCache) Set(name string, class, t uint16, record rr.RR, ttl uint32) error {
	node, err := c.Populate(name)
	if err != nil {
		return err
	}

	expire := time.Now().Add(time.Duration(ttl) * time.Second)
	node.SetData(class, t, record, expire)

	return nil
}

// Set sets (or adds) multiple new cache entry at (to) the same node
func (c *DefaultCache) SetMultiple(name string, class, t uint16, records []rr.RR, ttl uint32) error {
	node, err := c.Populate(name)
	if err != nil {
		return err
	}

	expire := time.Now().Add(time.Duration(ttl) * time.Second)
	for _, record := range records {
		node.SetData(class, t, record, expire)
	}

	return nil
}

// Package cache provides an in-memory cache to store answers from remote DNS servers
package cache

import (
	"errors"
	"time"

	"github.com/go-void/portal/pkg/tree"
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
	Lookup(string, uint16, uint16) (tree.Entry, Status, error)

	LookupQuestion(dns.Question) (tree.Entry, Status, error)

	Set(string, uint16, uint16, rr.RR, uint32) error

	SetMultiple(string, uint16, uint16, []rr.RR, uint32) error
}

// DefaultCache implements the Cache interface and stores
// RRs in an in-memory tree
type DefaultCache struct {
	*tree.Tree
}

// NewDefaultCache returns a new default in-memory tree cache
func NewDefaultCache() *DefaultCache {
	return &DefaultCache{
		Tree: tree.New(),
	}
}

// Lookup looks up a entry for name with class and type and returns the status and errors
// encountered along the way
func (c *DefaultCache) Lookup(name string, class, t uint16) (tree.Entry, Status, error) {
	node, err := c.Get(name)
	if err != nil {
		return tree.Entry{}, Miss, nil
	}

	entry, err := node.Entry(class, t)
	if err != nil {
		return tree.Entry{}, Miss, nil
	}

	status := Hit
	if entry.Expire.Before(time.Now()) {
		status = Expired
	}

	rr.UpdateTTL(entry.Record, entry.Expire)

	return tree.Entry{
		Record: entry.Record,
		Expire: entry.Expire,
	}, status, nil
}

// LookupQuestion is a convenience function to lookup a DNS question
func (c *DefaultCache) LookupQuestion(message dns.Question) (tree.Entry, Status, error) {
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

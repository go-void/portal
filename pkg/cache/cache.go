// Package cache provides an in-memory cache to store answers from remote DNS servers
package cache

import (
	"errors"

	"github.com/go-void/portal/pkg/logger"
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
	Lookup(string, uint16, uint16) ([]rr.RR, Status, error)

	LookupQuestion(dns.Question) ([]rr.RR, Status, error)

	Set(string, uint16, uint16, []rr.RR) error
}

// DefaultCache implements the Cache interface and stores
// RRs in an in-memory tree
type DefaultCache struct {
	*tree.Tree

	logger *logger.Logger
}

// NewDefaultCache returns a new default in-memory tree cache
func NewDefaultCache(l *logger.Logger) *DefaultCache {
	return &DefaultCache{
		Tree:   tree.New(),
		logger: l,
	}
}

// Lookup looks up a entry for name with class and type and returns the status and errors
// encountered along the way
func (c *DefaultCache) Lookup(name string, class, t uint16) ([]rr.RR, Status, error) {
	node, err := c.Get(name)
	if err != nil {
		return nil, Miss, nil
	}

	records, err := node.Entry(class, t)
	if err != nil {
		return nil, Miss, nil
	}

	// TODO (Techassi): Figure out how we should handle the expire of multiple RRs
	status := Hit
	// if entry.Expire.Before(time.Now()) {
	// 	status = Expired
	// }

	// rr.UpdateTTL(entry.Record, entry.Expire)

	return records, status, nil
}

// LookupQuestion is a convenience function to lookup a DNS question
func (c *DefaultCache) LookupQuestion(message dns.Question) ([]rr.RR, Status, error) {
	return c.Lookup(message.Name, message.Class, message.Type)
}

// Set sets (or adds) a new cache entry
func (c *DefaultCache) Set(name string, class, t uint16, records []rr.RR) error {
	node, err := c.Populate(name)
	if err != nil {
		return err
	}

	node.SetEntry(class, t, records)
	return nil
}

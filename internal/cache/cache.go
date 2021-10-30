// Package cache provides an in-memory cache to store answers from remote DNS servers
package cache

import (
	"time"

	"github.com/go-void/portal/internal/tree"
	"github.com/go-void/portal/internal/types/dns"
	"github.com/go-void/portal/internal/types/rr"
)

type Entry struct {
	Data   rr.RR
	Expire time.Time
}

type Cache interface {
	Get(string, uint16, uint16) (Entry, Status, error)

	GetQuestion(dns.Question) (Entry, Status, error)

	Set(string, uint16, uint16, rr.RR, uint32) error
}

type DefaultCache struct {
	Tree *tree.Tree
}

func NewDefaultCache() *DefaultCache {
	return &DefaultCache{
		Tree: tree.New(),
	}
}

func (c *DefaultCache) Get(name string, class, t uint16) (Entry, Status, error) {
	node, err := c.Tree.Get(name)
	if err != nil {
		return Entry{}, Miss, nil
	}

	nodeRecord, err := node.Record(class, t)
	if err != nil {
		return Entry{}, Miss, err
	}

	record := nodeRecord.RR
	record.SetHeader(rr.Header{
		Name:     name,
		Class:    class,
		Type:     t,
		TTL:      3600,
		RDLength: record.Len(),
	})

	return Entry{
		Data:   record,
		Expire: nodeRecord.Expire,
	}, Hit, nil
}

func (c *DefaultCache) GetQuestion(message dns.Question) (Entry, Status, error) {
	return c.Get(message.Name, message.Class, message.Type)
}

func (c *DefaultCache) Set(name string, class, t uint16, record rr.RR, ttl uint32) error {
	node, err := c.Tree.Populate(name)
	if err != nil {
		return err
	}

	node.SetData(class, t, record, ttl)
	return nil
}

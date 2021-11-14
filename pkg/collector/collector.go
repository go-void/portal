// Package collector provides functions to collect statistics about queries, filters and more
package collector

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/go-void/portal/pkg/config"
	"github.com/go-void/portal/pkg/types/dns"
	"github.com/go-void/portal/pkg/types/rr"
)

// Collector describes an interface which is able to collect
// stats about queries, filters and more. The collector first
// stores this data in memory. At every persistence interval
// this data gets persisted via the CollectorStore interface
type Collector interface {
	// GetEntries returns a slice of query stats which
	// the collector currently has stored in memory
	GetEntries() ([]Entry, error)

	// AddEntry adds a new query entry to the
	// collector
	AddEntry(Entry) error

	// FlushEntries flushes in-memory entries
	// via the store interface (to disk)
	FlushEntries() error

	// Run runs the flush interval
	Run()
}

type DefaultCollector struct {
	Store CollectorStore

	Interval   time.Duration
	LastFlush  time.Time
	MaxEntries int
	Anonymize  bool
	Enabled    bool

	entries []Entry
	lock    sync.Mutex
}

func NewEntry(question dns.Question, answer rr.RR, queryTime time.Duration, ip net.IP) Entry {
	return Entry{
		QueryTime: queryTime,
		Question:  question,
		Answer:    answer,
		Result:    "",
		ClientIP:  ip,
		Filtered:  false,
		Cached:    false,
	}
}

func NewCachedEntry(question dns.Question, answer rr.RR, queryTime time.Duration, ip net.IP) Entry {
	entry := NewEntry(question, answer, queryTime, ip)
	entry.Cached = true
	return entry
}

func NewFilteredEntry(question dns.Question, answer rr.RR, queryTime time.Duration, ip net.IP) Entry {
	entry := NewEntry(question, answer, queryTime, ip)
	entry.Filtered = true
	return entry
}

func NewDefault(opt config.CollectorOptions) *DefaultCollector {
	interval := time.Duration(opt.Interval) * time.Second

	return &DefaultCollector{
		Store:      NewDefaultStore(),
		Interval:   interval,
		MaxEntries: opt.MaxEntries,
		Anonymize:  opt.Anonymize,
		Enabled:    opt.Enabled,
	}
}

func (c *DefaultCollector) GetEntries() ([]Entry, error) {
	return c.entries, nil
}

func (c *DefaultCollector) AddEntry(entry Entry) error {
	c.entries = append(c.entries, entry)
	fmt.Println(c.entries)
	if len(c.entries) >= c.MaxEntries {
		return c.FlushEntries()
	}
	return nil
}

func (c *DefaultCollector) FlushEntries() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if len(c.entries) == 0 {
		return nil
	}

	c.LastFlush = time.Now()
	c.entries = []Entry{}

	return nil
}

func (c *DefaultCollector) Run() {
	qt := time.NewTicker(c.Interval)

	go func() {
		for range qt.C {
			if time.Since(c.LastFlush) < (c.Interval / 2) {
				continue
			}

			err := c.FlushEntries()
			if err != nil {
				// Handle
				fmt.Println(err)
			}
		}
	}()
}

package cache

import (
	"time"

	"github.com/go-void/portal/internal/types/rr"
)

type Node struct {
	parent   *Node
	children map[string]Node
	entries  map[uint16]Entry
}

type Entry struct {
	Record rr.RR
	Expire time.Time
}

// Parent returns this node's parent node
func (n *Node) Parent() *Node {
	return n.parent
}

// Child returns this node's child identified by 'name" or an error if this child doesn't exist
func (n *Node) Child(name string) (Node, error) {
	if node, ok := n.children[name]; ok {
		return node, nil
	}
	return Node{}, ErrNodeNotFound
}

// AddChild adds a child to this node or returns an error if the child already exists
func (n *Node) AddChild(name string, child Node) error {
	if _, ok := n.children[name]; ok {
		return ErrChildAlreadyExists
	}
	n.children[name] = child
	return nil
}

// Record returns a stored record with class and type
func (n *Node) Entry(class, t uint16) (Entry, error) {
	if record, ok := n.entries[class*100+t]; ok {
		return record, nil
	}
	return Entry{}, ErrNoSuchData
}

// SetData sets data for type t
func (n *Node) SetData(class, t uint16, record rr.RR, expire time.Time) {
	n.entries[class*100+t] = Entry{
		Expire: expire,
		Record: record,
	}
}

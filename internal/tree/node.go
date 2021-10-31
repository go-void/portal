package tree

import (
	"time"

	"github.com/go-void/portal/internal/types/rr"
)

// Node describes a single node which keeps tracks of it's children,
// data and parent node
type Node struct {
	parent   *Node
	children map[string]Node
	data     map[uint16]Record
}

type Record struct {
	Expire time.Time
	RR     rr.RR
}

// Parent returns this node's parent node
func (n *Node) Parent() *Node {
	return n.parent
}

// Child returns this node's child identified by 'name" or
// an error if this child doesn't exist
func (n *Node) Child(name string) (Node, error) {
	if node, ok := n.children[name]; ok {
		return node, nil
	}
	return Node{}, ErrNodeNotFound
}

// AddChild adds a child to this node or returns an error
// if the child already exists
func (n *Node) AddChild(name string, child Node) error {
	if _, ok := n.children[name]; ok {
		return ErrChildAlreadyExists
	}
	n.children[name] = child
	return nil
}

// Record returns a stored record with class and type
func (n *Node) Record(class, t uint16) (Record, error) {
	if record, ok := n.data[class*100+t]; ok {
		return record, nil
	}
	return Record{}, ErrNoSuchData
}

// SetData sets data for type t
func (n *Node) SetData(class, t uint16, record rr.RR, expire time.Time) {
	n.data[class*100+t] = Record{
		Expire: expire,
		RR:     record,
	}
}

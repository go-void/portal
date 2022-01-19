package tree

import (
	"github.com/go-void/portal/pkg/types/rr"
)

type Node struct {
	parent   *Node
	children map[string]Node
	records  map[uint16][]rr.RR
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
func (n *Node) Records(class, t uint16) ([]rr.RR, error) {
	if entry, ok := n.records[class*100+t]; ok {
		return entry, nil
	}
	return nil, ErrNoSuchData
}

// AddRecords adds records to this node
func (n *Node) AddRecords(records []rr.RR) {
	isSame := false

	for i := 0; i < len(records); i++ {
		key := records[i].Header().Class*100 + records[i].Header().Type

		for j := 0; j < len(n.records[key]); j++ {
			if records[i].IsSame(n.records[key][j]) {
				isSame = true
				break
			}
		}

		if isSame {
			continue
		}

		n.records[key] = append(n.records[key], records[i])
	}
}

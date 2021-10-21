// Package tree implements a in-memory tree / node structure to store
// and retrieve resource record data
package tree

import (
	"errors"

	"github.com/go-void/portal/internal/utils"
)

var (
	ErrNodeNotFound       = errors.New("node not found in tree")
	ErrNoSuchData         = errors.New("no such data")
	ErrChildAlreadyExists = errors.New("child already exists")
)

// Tree describes a tree structure which stores data for DNS labels
type Tree struct {
	root Node
}

func New() *Tree {
	return &Tree{
		root: Node{
			parent:   nil,
			children: make(map[string]Node),
			data:     make(map[uint16]interface{}),
		},
	}
}

// Get retrieves a node via a name. Example: example.com traverses
// the tree like . -> com -> example
func (t *Tree) Get(name string) (Node, error) {
	nodes, err := t.Walk(name)
	if err != nil {
		return Node{}, err
	}
	return nodes[len(nodes)-1], nil
}

// Walk traverses the tree until the end of labels is reached which
// returns a list of node or a non-existent node is encountered and
// an error is returned
func (t *Tree) Walk(name string) ([]Node, error) {
	var (
		current = t.root
		nodes   = []Node{}
		names   = utils.LabelsFromRoot(name)
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

// Populate traverses the tree and adds nodes along the path which
// don't exist yet
func (t *Tree) Populate(name string) (Node, error) {
	var current = t.root
	var names = utils.LabelsFromRoot(name)

	for _, name := range names {
		if name == "" || name == "." {
			continue
		}

		node, err := current.Child(name)
		if err != nil {
			node := Node{
				parent:   &current,
				children: make(map[string]Node),
				data:     make(map[uint16]interface{}),
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

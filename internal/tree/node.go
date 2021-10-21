package tree

// Node describes a single node which keeps tracks of it's children,
// data and parent node
type Node struct {
	parent   *Node
	children map[string]Node
	data     map[uint16]interface{}
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

// Data returns stored data for type t
func (n *Node) Data(class, t uint16) (interface{}, error) {
	if data, ok := n.data[class*100+t]; ok {
		return data, nil
	}
	return nil, ErrNoSuchData
}

// SetData sets data for type t
func (n *Node) SetData(class, t uint16, data interface{}) {
	n.data[class*100+t] = data
}

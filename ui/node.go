package ui

import (
	"image"
)

type Node struct {
	parent     *Node
	firstChild *Node
	lastChild  *Node
	next, prev *Node
	bounds     image.Rectangle
}

func (n *Node) init(parent *Node) {
	n.parent = parent
}

func (n *Node) node() *Node {
	return n
}

// Append appends child to the node as the new LastChild.
func (n *Node) Append(child *Node) {
	if child.parent != nil {
		child.unlinkParent()
	}
	if n.lastChild == nil {
		n.firstChild = child
		n.lastChild = child
	} else {
		child.prev = n.lastChild
		n.lastChild.next = child
		n.lastChild = child
	}
	child.init(n)
}

// Parent returns the parent node.
func (n *Node) Parent() *Node {
	return n.parent
}

// FirstChild returns the first child node.
func (n *Node) FirstChild() *Node {
	return n.firstChild
}

// LastChild returns the last child node.
func (n *Node) LastChild() *Node {
	return n.lastChild
}

// Next returns the next sibling node.
func (n *Node) Next() *Node {
	return n.next
}

// Prev returns the previous sibling node.
func (n *Node) Prev() *Node {
	return n.prev
}

// Delete unlinks the node and releases its resources.
func (n *Node) Delete() {
	n.unlinkParent()
	// TODO: delete children resources
	n.parent = nil
}

// Bounds returns the visual boundary of the node.
func (n *Node) Bounds() image.Rectangle {
	return n.bounds
}

// SetBounds changes the visual boundary of the node.
func (n *Node) SetBounds(r image.Rectangle) {
	n.bounds = r
}

// Within returns true if pt is within the visual boundary of the node.
func (n *Node) Within(pt image.Point) bool {
	switch {
	case pt.X < n.bounds.Min.X,
		pt.Y < n.bounds.Min.Y,
		pt.X > n.bounds.Max.X,
		pt.Y > n.bounds.Max.Y:
		return false
	default:
		return true
	}
}

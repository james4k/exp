package ui

import (
	"image"
)

// Element is the basis of a Component node which can take user input.
type Element struct {
	Node
	self     Component
	children []Component
}

func (e *Element) initElement(cmpt Component) {
	e.self = cmpt
}

func (e *Element) Focus(pt image.Point) Component {
	for _, cmpt := range e.children {
		n := cmpt.node()
		if n.Within(pt) {
			cmpt = cmpt.Focus(pt)
			if cmpt != nil {
				return cmpt
			}
		}
	}
	return e.self
}

func (e *Element) Children() Components {
	return Components{cc: e.children}
}

type Component interface {
	init(*Node)
	node() *Node
	initElement(Component)

	Enter()
	Exit()
	Children(pt image.Point) Components
	Focus(pt image.Point) Component
	Mouse(Mouse)
	Keyboard(Keyboard)
}

type Components struct {
	cc []Component
	i  int
}

func (c *Components) Next() bool {
	if c.i < len(c.cc) {
		c.i++
		return true
	}
	return false
}

func (c *Components) Component() Component {
	return c.cc[c.i-1]
}

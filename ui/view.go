package ui

import "image"

// View is composed of a Box for spatial properties, and usually
// additional fields to describe state and visuals.
type View interface {
	box() *Box
	Bounds() image.Rectangle
	Subviews() int
	Sub(i int) View
}

// Component is a View which responds to events and manages subviews.
type Component interface {
	View
	Run(*Controller)
}

// Box describes the spatial and hierarchical properties of a View.
type Box struct {
	kids   []View
	bounds image.Rectangle
	ctl    *Controller
}

// setup initializes a View and its Box. Boots up Components.
func setup(parent *Box, bounds image.Rectangle, view View) {
	box := view.box()
	if box == nil {
		panic("ui: box must be non-nil")
	}
	*box = Box{
		bounds: bounds,
	}
	if comp, ok := view.(Component); ok {
		ctl := &Controller{
			box: box,
		}
		if parent == nil {
			ctl.recv = make(chan message, 1)
			ctl.chain = make(chan message, 1)
		} else {
			ctl.recv = make(chan message, 1)
			ctl.chain = parent.ctl.recv
		}
		box.ctl = ctl
		go comp.Run(ctl)
	}
}

func (b *Box) box() *Box {
	return b
}

func (b *Box) Subviews() int {
	return len(b.kids)
}

func (b *Box) Sub(i int) View {
	return b.kids[i]
}

func (b *Box) Size() (w, h int) {
	w = b.bounds.Dx()
	h = b.bounds.Dy()
	return
}

func (b *Box) Bounds() image.Rectangle {
	return b.bounds
}

func (b *Box) SetBounds(rect image.Rectangle) {
	b.bounds = rect
}

func (b *Box) hitTest(pt image.Point) *Box {
	if pt.In(b.bounds) {
		for _, k := range b.kids {
			target := k.box().hitTest(pt)
			if target != nil {
				return target
			}
		}
		return b
	}
	return nil
}

func (b *Box) close() {
	if b.ctl != nil && b.ctl.recv != nil {
		close(b.ctl.recv)
	}
	for _, k := range b.kids {
		k.box().close()
	}
}

// message carries an event value
type message struct {
	event interface{}
	done  chan<- message
}

// Controller controls a view and its subviews, and listens for events.
type Controller struct {
	box   *Box
	recv  chan message
	chain chan message
	msg   message
}

func (c *Controller) Listen() (event interface{}, ok bool) {
	if c.msg.done != nil {
		c.msg.done <- c.msg
	}
	msg, ok := <-c.recv
	if !ok {
		c.chain = nil
		c.msg.done = nil
		return nil, false
	}
	event = msg.event
	c.msg = msg
	return
}

func (c *Controller) Append(subviews ...View) {
	c.box.kids = append(c.box.kids, subviews...)
	for _, v := range subviews {
		setup(c.box, c.box.Bounds(), v)
	}
}

func (c *Controller) Remove(subviews ...View) {
	kids := c.box.kids
	for _, v := range subviews {
		for i, k := range kids {
			if k == v {
				copy(kids[i:], kids[i+1:])
				kids = kids[:len(kids)-1]
				break
			}
		}
	}
}

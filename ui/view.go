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

// Component is a View which receives events and manages subviews.
type Component interface {
	View
	Receiver
}

// Receiver receives events.
type Receiver interface {
	Receive(ctl *Controller, event interface{})
}

// Box describes the spatial and hierarchical properties of a View.
type Box struct {
	kids   []View
	bounds image.Rectangle
	ctl    *Controller
}

// setup initializes a View and its Box. Mounts Components.
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
			box:  box,
			comp: comp,
		}
		box.ctl = ctl
		comp.Receive(ctl, Mount{})
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

func (b *Box) send(event interface{}) {
	if b.ctl != nil {
		b.ctl.comp.Receive(b.ctl, event)
	}
}

func (b *Box) unmount() {
	for _, k := range b.kids {
		k.box().unmount()
	}
	b.send(Unmount{})
}

// Controller controls a view and its subviews.
type Controller struct {
	box  *Box
	comp Component
}

func (c *Controller) Mount(subviews ...View) {
	c.box.kids = append(c.box.kids, subviews...)
	for _, v := range subviews {
		setup(c.box, c.box.Bounds(), v)
	}
}

func (c *Controller) Unmount(subviews ...View) {
	kids := c.box.kids
	for _, v := range subviews {
		for i, k := range kids {
			if k == v {
				k.box().unmount()
				copy(kids[i:], kids[i+1:])
				kids = kids[:len(kids)-1]
				break
			}
		}
	}
}

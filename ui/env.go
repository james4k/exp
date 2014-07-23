package ui

import "image"

// Environment
type Environment interface {
	Size() (width, height int, pixelRatio float32)
	Listener
}

// Listener listens for incoming events.
type Listener interface {
	// Listen blocks until an event is received, or there are no more
	// events. If called once, must be called again until ok returns
	// false.
	Listen() (event interface{}, ok bool)
}

// ListenFor listens for an event of a specific type. Returns true when
// an event was recieved of the type pointed to by p, and returns false
// when there are no more events to listen for. Events of all other
// types are ignored.
func ListenFor(p interface{}, l Listener) bool {
	return false
}

// CharTyped given when a Unicode character is input. Meant for text
// entry, not for physical key state (pressed/released), modifier keys,
// etc.
type CharTyped struct {
	Char rune
}

// MouseUpdate given on changes to the position of the cursor, or the
// state of the mouse buttons.
type MouseUpdate struct {
	image.Point
	Left, Right bool
}

type MouseEnter struct {
}

type MouseLeave struct {
}

// SizeUpdate given on window resize.
type SizeUpdate struct {
	Width  int
	Height int
}

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
	// events, ie. when ok is false.
	Listen() (event interface{}, ok bool)
}

type Mount struct {
}

type Unmount struct {
}

type MouseState struct {
	image.Point
	Left, Right bool
}

// MouseUpdate given on changes to the position of the cursor, or the
// state of the mouse buttons.
type MouseUpdate struct {
	MouseState
	Previous MouseState
}

type MouseEnter struct {
}

type MouseLeave struct {
}

type FocusGained struct {
}

type FocusLost struct {
}

type KeyboardUpdate struct {
}

type KeyDown struct {
	Key
}
type KeyUp struct {
	Key
}
type KeyRepeat struct {
	Key
}

// UnicodeTyped given when a Unicode character is input. Meant for text
// entry, not for physical key state (pressed/released), modifier keys,
// etc.
type UnicodeTyped struct {
	C rune
}

// SizeUpdate given on window resize.
type SizeUpdate struct {
	Width  int
	Height int
}

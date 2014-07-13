package ui

import (
	"image"

	"j4k.co/exp/ui/graphics"
)

// Environment represents a window's environment.
type Environment interface {
	Keyboard() <-chan Keyboard
	Mouse() <-chan Mouse
	View() <-chan View
	Graphics() graphics.Context
}

// Keyboard is just a character, for now. This may represent the entire
// keyboard state in the future, or perhaps a separate interface.
type Keyboard struct {
	Typed rune
}

// Mouse is the position of the cursor, and the state of the
// mouse buttons.
type Mouse struct {
	image.Point
	Left,
	Right bool
}

// View is the user's view of the UI.
type View struct {
	Width  int
	Height int

	readyc chan struct{}
}

// Ready should be called when the surface is ready after the view
// change.
func (v View) Ready() {
	v.readyc <- struct{}{}
}

package ui

// Environment represents a window's environment.
type Environment interface {
	Keyboard() <-chan Keyboard
	Mouse() <-chan Mouse
	View() <-chan View
	Graphics() Graphics
}

// Keyboard is just a character, for now. This may represent the entire
// keyboard state in the future, or perhaps a separate interface.
type Keyboard struct {
	Typed rune
}

// Mouse is the position of the cursor, and the state of the
// mouse buttons.
type Mouse struct {
	X, Y int
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

// Graphics represents the GL context to use.
type Graphics interface {
	// Lock prepares the GL context and thread for GL calls.
	Lock()
	// Unlock releases the GL thread for the goroutine.
	Unlock()
	// SwapBuffers swaps the front/back buffers, waiting on vsync.
	SwapBuffers()
}

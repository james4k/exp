package graphics

// Context represents a GL context for use by the graphics system.
type Context interface {
	// Lock prepares the GL context and thread for GL calls.
	Lock()
	// Unlock releases the GL thread for the goroutine.
	Unlock()
	// SwapBuffers swaps the front/back buffers, waiting on vsync.
	SwapBuffers()
}

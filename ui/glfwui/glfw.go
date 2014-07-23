// Package glfwui provides ui.Environment creation via GLFW, a
// cross-platform library for creating windows, OpenGL contexts, and
// managing input and events.
package glfwui

import (
	"image"
	"runtime"

	"j4k.co/ui"
	//"j4k.co/exp/ui/graphics"

	"github.com/go-gl/glfw3"
)

type Window struct {
	w      *glfw3.Window
	eventc chan interface{}
	waitc  chan struct{}

	mouse ui.MouseUpdate
}

// Open opens a new OS window via GLFW.
func Open(width, height int, title string) (*Window, error) {
	var wnd *glfw3.Window
	errc := make(chan error)
	mainc <- func() {
		var err error
		wnd, err = glfw3.CreateWindow(width, height, title, nil, nil)
		errc <- err
		if err == nil {
			wnd.SetInputMode(glfw3.Cursor, glfw3.CursorNormal)
			wnd.Restore()
		}
	}
	err := <-errc
	if err != nil {
		return nil, err
	}
	return openFromWindow(wnd)
}

func openFromWindow(wnd *glfw3.Window) (*Window, error) {
	w := &Window{
		w: wnd,
	}
	w.init()
	return w, nil
}

func (w *Window) GLFWWindow() *glfw3.Window {
	return w.w
}

func (w *Window) init() {
	// note: eventc must be unbuffered for the Listen/waitc setup to
	// work
	w.eventc = make(chan interface{})
	w.w.SetCharacterCallback(w.onCharPress)
	w.w.SetMouseButtonCallback(w.onMouseButton)
	w.w.SetCursorPositionCallback(w.onCursorPos)
	w.w.SetSizeCallback(w.onResize)
	w.w.SetCloseCallback(w.onClose)
}

func (w *Window) MakeContextCurrent() {
	runtime.LockOSThread()
	w.w.MakeContextCurrent()
}

func (w *Window) DetachContext() {
	glfw3.DetachCurrentContext()
	runtime.UnlockOSThread()
}

func (w *Window) SwapBuffers() {
	w.w.SwapBuffers()
}

func (w *Window) Size() (ww, h int, pixelRatio float32) {
	_, fh := w.w.GetFramebufferSize()
	ww, h = w.w.GetSize()
	pixelRatio = float32(fh) / float32(h)
	return
}

func (w Window) dispatch(event interface{}) {
	w.eventc <- event
	<-w.waitc
}

func (w *Window) close() {
	close(w.eventc)
}

func (w *Window) Listen() (event interface{}, ok bool) {
	if w.waitc == nil {
		w.waitc = make(chan struct{})
	} else {
		w.waitc <- struct{}{}
	}
	e, ok := <-w.eventc
	if !ok {
		w.w.Destroy()
		// TODO: this breaks support for multiple windows
		close(mainc)
		return nil, false
	}
	return e, true
}

func (w *Window) onCharPress(wnd *glfw3.Window, char uint) {
	w.dispatch(ui.CharTyped{
		Char: rune(char),
	})
}

func (w *Window) onMouseButton(wnd *glfw3.Window, btn glfw3.MouseButton, action glfw3.Action, mod glfw3.ModifierKey) {
	switch btn {
	case glfw3.MouseButton1:
		w.mouse.Left = (action != glfw3.Release)
		w.dispatch(w.mouse)
	}
}

func (w *Window) onCursorPos(wnd *glfw3.Window, x, y float64) {
	w.mouse.Point = image.Pt(int(x), int(y))
	w.dispatch(w.mouse)
}

func (w *Window) onResize(wnd *glfw3.Window, ww, h int) {
	w.dispatch(ui.SizeUpdate{
		Width:  ww,
		Height: h,
	})
}

func (w *Window) onClose(wnd *glfw3.Window) {
	w.close()
}

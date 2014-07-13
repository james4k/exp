// Package glfwenv provides ui.Environment creation via GLFW, a
// cross-platform library for creating windows, OpenGL contexts and
// managing input and events.
package glfwenv

import (
	"image"
	"runtime"
	"sync"

	"j4k.co/exp/ui"
	"j4k.co/exp/ui/graphics"

	glfw "github.com/go-gl/glfw3"
)

type context struct {
	w  *glfw.Window
	mu sync.Mutex
}

func (c *context) Lock() {
	runtime.LockOSThread()
	c.mu.Lock()
	c.w.MakeContextCurrent()
}

func (c *context) Unlock() {
	c.mu.Unlock()
	runtime.UnlockOSThread()
}

func (c *context) SwapBuffers() {
	c.w.SwapBuffers()
}

type env struct {
	*glfw.Window
	k chan ui.Keyboard
	m chan ui.Mouse
	v chan ui.View
	g graphics.Context

	mouse ui.Mouse
}

// New creates a new OS window via GLFW as a ui.Environment.
func New(width, height int, title string) (ui.Environment, error) {
	var wnd *glfw.Window
	errc := make(chan error)
	mainc <- func() {
		var err error
		wnd, err = glfw.CreateWindow(width, height, title, nil, nil)
		errc <- err
		if err == nil {
			wnd.SetInputMode(glfw.Cursor, glfw.CursorNormal)
		}
	}
	err := <-errc
	if err != nil {
		return nil, err
	}
	w := &env{
		Window: wnd,
	}
	w.init()
	return w, nil
}

func (g *env) init() {
	g.k = make(chan ui.Keyboard)
	g.m = make(chan ui.Mouse)
	g.v = make(chan ui.View)
	g.g = &context{
		w: g.Window,
	}
	g.SetCharacterCallback(g.onCharPress)
	g.SetMouseButtonCallback(g.onMouseButton)
	g.SetCursorPositionCallback(g.onCursorPos)
	g.SetSizeCallback(g.onResize)
}

func (g *env) Keyboard() <-chan ui.Keyboard {
	return g.k
}

func (g *env) Mouse() <-chan ui.Mouse {
	return g.m
}

func (g *env) View() <-chan ui.View {
	return g.v
}

func (g *env) Graphics() graphics.Context {
	return g.g
}

func (g *env) onCharPress(wnd *glfw.Window, char uint) {
	g.k <- ui.Keyboard{
		Typed: rune(char),
	}
}

func (g *env) onMouseButton(wnd *glfw.Window, btn glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	switch btn {
	case glfw.MouseButton1:
		g.mouse.Left = (action != glfw.Release)
		g.m <- g.mouse
	}
}

func (g *env) onCursorPos(wnd *glfw.Window, x, y float64) {
	g.mouse.Point = image.Pt(int(x), int(y))
	g.m <- g.mouse
}

func (g *env) onResize(wnd *glfw.Window, w, h int) {
	g.v <- ui.View{
		Width:  w,
		Height: h,
	}
}

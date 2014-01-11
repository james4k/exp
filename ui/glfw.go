package ui

import (
	glfw "github.com/go-gl/glfw3"
	"runtime"
	"sync"
)

type glfwGraphics struct {
	w  *glfw.Window
	mu sync.Mutex
}

func (g *glfwGraphics) Lock() {
	runtime.LockOSThread()
	g.mu.Lock()
	g.w.MakeContextCurrent()
}

func (g *glfwGraphics) Unlock() {
	g.mu.Unlock()
	runtime.UnlockOSThread()
}

func (g *glfwGraphics) SwapBuffers() {
	g.w.SwapBuffers()
}

type glfwEnv struct {
	*glfw.Window
	keyboard chan Keyboard
	mouse    chan Mouse
	view     chan View
	graphics Graphics
}

func glfwEnvironment(width, height int, title string) (*glfwEnv, error) {
	var wnd *glfw.Window
	errc := make(chan error)
	mainc <- func() {
		var err error
		wnd, err = glfw.CreateWindow(width, height, title, nil, nil)
		errc <- err
	}
	err := <-errc
	if err != nil {
		return nil, err
	}
	w := &glfwEnv{
		Window: wnd,
	}
	w.init()
	return w, nil
}

func (g *glfwEnv) init() {
	g.keyboard = make(chan Keyboard)
	g.mouse = make(chan Mouse)
	g.view = make(chan View)
	g.SetCharacterCallback(g.onCharPress)
	g.graphics = &glfwGraphics{
		w: g.Window,
	}
}

func (g *glfwEnv) Keyboard() <-chan Keyboard {
	return g.keyboard
}

func (g *glfwEnv) Mouse() <-chan Mouse {
	return g.mouse
}

func (g *glfwEnv) View() <-chan View {
	return g.view
}

func (g *glfwEnv) Graphics() Graphics {
	return g.graphics
}

func (g *glfwEnv) onCharPress(wnd *glfw.Window, char uint) {
	g.keyboard <- Keyboard{
		Typed: rune(char),
	}
}

func (g *glfwEnv) onResize(wnd *glfw.Window, w, h int) {
	g.view <- View{
		Width:  w,
		Height: h,
	}
}

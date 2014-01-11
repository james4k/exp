package ui

import (
	"image"
	//"io"
	//"log"
	"sync"
)

type void struct{}

// Window is the root of the visual hierarchy in the user interface.
type Window struct {
	nodes   []*Node
	body    Element
	wet     chan void
	rlist   []renderable
	rlistmu sync.Mutex
}

// Create creates a new OS window.
func Create(width, height int, title string) (*Window, error) {
	wnd, err := glfwEnvironment(width, height, title)
	if err != nil {
		return nil, err
	}
	w := &Window{
		wet: make(chan void, 1),
	}
	go w.Run(wnd)
	return w, nil
}

// Run executes the window using env.
func (w *Window) Run(env Environment) {
	keybd := env.Keyboard()
	mouse := env.Mouse()
	view := env.View()
	graphics := env.Graphics()

	w.init()
	for {
		select {
		case k := <-keybd:
			w.keyboard(k)
		case m := <-mouse:
			w.mouse(m)
		case v := <-view:
			w.view(v, graphics)
		case <-w.wet:
			w.paint(graphics)
		}
	}
}

func (w *Window) Component(parent *Node, c Component) {
	// TODO: panic if component already initialized
	parent.Append(c.node())
	c.Enter()
}

func (w *Window) init() {
	//w.nodes = append(w.nodes, w.body)
}

// mouse updates mouse focus and passes mouse input to the focused
// element.
func (w *Window) mouse(m Mouse) {
	focus := w.body.Focus(image.Pt(m.X, m.Y))
	if focus != nil {
		focus.Mouse(m)
	}
}

// keyboard sends keyboard input to the keyboard focused element.
func (w *Window) keyboard(k Keyboard) {
	//if w.keyfocus != nil {
	//w.keyfocus.Keyboard(k)
	//}
}

func (w *Window) view(v View, g Graphics) {
	// TODO: server should just wait for swapbuffers
	defer v.Ready()
	w.body.SetBounds(image.Rect(0, 0, v.Width, v.Height))
	w.paint(g)
}

func (w *Window) paint(g Graphics) {
	//w.Layout()
	w.render(g)
}

func (w *Window) render(g Graphics) {
	w.prepareRenderList()
	g.Lock()
	defer g.Unlock()
	w.renderInit()
	w.renderUpload()
	w.renderList()
	g.SwapBuffers()
}

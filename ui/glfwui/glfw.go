// Package glfwui provides ui.Environment creation via GLFW, a
// cross-platform library for creating windows, OpenGL contexts, and
// managing input and events.
package glfwui

import (
	"image"
	"runtime"

	"j4k.co/exp/ui"

	"github.com/go-gl/glfw3"
)

type Window struct {
	w           *glfw3.Window
	eventc      chan interface{}
	waitc       chan struct{}
	haslistened bool

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
	w.waitc = make(chan struct{})
	w.w.SetCharacterCallback(w.onCharPress)
	w.w.SetKeyCallback(w.onKeyPress)
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

func (w *Window) dispatch(event interface{}) {
	w.eventc <- event
	<-w.waitc
}

func (w *Window) close() {
	close(w.eventc)
}

func (w *Window) Listen() (event interface{}, ok bool) {
	if !w.haslistened {
		w.haslistened = true
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
	// ignore OS X arrow keys, etc... glfw will get IME support in 3.2
	// and so input will likely look a lot different eventually,
	// anyways.
	if char >= 0xf700 {
		return
	}
	w.dispatch(ui.UnicodeTyped{
		C: rune(char),
	})
}

func (w *Window) onKeyPress(wnd *glfw3.Window, key glfw3.Key, scancode int, action glfw3.Action, mod glfw3.ModifierKey) {
	var s []byte
	if mod&glfw3.ModShift != 0 {
		s = append(s, '$')
	}
	if mod&glfw3.ModControl != 0 {
		s = append(s, '^')
	}
	if mod&glfw3.ModAlt != 0 {
		s = append(s, '~')
	}
	// TODO: super?
	s = translateKey(s, key)
	switch action {
	case glfw3.Press:
		w.dispatch(ui.KeyDown{
			Key: ui.Key(s),
		})
	case glfw3.Release:
		w.dispatch(ui.KeyUp{
			Key: ui.Key(s),
		})
	case glfw3.Repeat:
		w.dispatch(ui.KeyRepeat{
			Key: ui.Key(s),
		})
	}
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

var keymap = map[glfw3.Key][]byte{
	glfw3.KeySpace:      {' '},
	glfw3.KeyApostrophe: {'\''},
	glfw3.KeyComma:      {','},
	glfw3.KeyMinus:      {'-'},
	glfw3.KeyPeriod:     {'.'},
	glfw3.KeySlash:      {'/'},

	glfw3.KeySemicolon:    {';'},
	glfw3.KeyEqual:        {'='},
	glfw3.KeyLeftBracket:  {'['},
	glfw3.KeyRightBracket: {']'},
	glfw3.KeyBackslash:    {'\\'},
	glfw3.KeyGraveAccent:  {'`'},

	glfw3.KeyEscape:    []byte(ui.Escape),
	glfw3.KeyEnter:     []byte(ui.Enter),
	glfw3.KeyTab:       []byte(ui.Tab),
	glfw3.KeyBackspace: []byte(ui.Backspace),
	glfw3.KeyInsert:    []byte(ui.Insert),
	glfw3.KeyDelete:    []byte(ui.Delete),
	glfw3.KeyLeft:      []byte(ui.Left),
	glfw3.KeyUp:        []byte(ui.Up),
	glfw3.KeyRight:     []byte(ui.Right),
	glfw3.KeyDown:      []byte(ui.Down),

	glfw3.Key0: {'0'},
	glfw3.Key1: {'1'},
	glfw3.Key2: {'2'},
	glfw3.Key3: {'3'},
	glfw3.Key4: {'4'},
	glfw3.Key5: {'5'},
	glfw3.Key6: {'6'},
	glfw3.Key7: {'7'},
	glfw3.Key8: {'8'},
	glfw3.Key9: {'9'},

	glfw3.KeyKp0:        {'#', '0'},
	glfw3.KeyKp1:        {'#', '1'},
	glfw3.KeyKp2:        {'#', '2'},
	glfw3.KeyKp3:        {'#', '3'},
	glfw3.KeyKp4:        {'#', '4'},
	glfw3.KeyKp5:        {'#', '5'},
	glfw3.KeyKp6:        {'#', '6'},
	glfw3.KeyKp7:        {'#', '7'},
	glfw3.KeyKp8:        {'#', '8'},
	glfw3.KeyKp9:        {'#', '9'},
	glfw3.KeyKpDecimal:  {'#', '.'},
	glfw3.KeyKpDivide:   {'#', '/'},
	glfw3.KeyKpMultiply: {'#', '*'},
	glfw3.KeyKpSubtract: {'#', '-'},
	glfw3.KeyKpAdd:      {'#', '+'},
	glfw3.KeyKpEnter:    {'#', 'e'},
	glfw3.KeyKpEqual:    {'#', '='},

	glfw3.KeyA: {'a'},
	glfw3.KeyB: {'b'},
	glfw3.KeyC: {'c'},
	glfw3.KeyD: {'d'},
	glfw3.KeyE: {'e'},
	glfw3.KeyF: {'f'},
	glfw3.KeyG: {'g'},
	glfw3.KeyH: {'h'},
	glfw3.KeyI: {'i'},
	glfw3.KeyJ: {'j'},
	glfw3.KeyK: {'k'},
	glfw3.KeyL: {'l'},
	glfw3.KeyM: {'m'},
	glfw3.KeyN: {'n'},
	glfw3.KeyO: {'o'},
	glfw3.KeyP: {'p'},
	glfw3.KeyQ: {'q'},
	glfw3.KeyR: {'r'},
	glfw3.KeyS: {'s'},
	glfw3.KeyT: {'t'},
	glfw3.KeyU: {'u'},
	glfw3.KeyV: {'v'},
	glfw3.KeyW: {'w'},
	glfw3.KeyX: {'x'},
	glfw3.KeyY: {'y'},
	glfw3.KeyZ: {'z'},

	glfw3.KeyF1:  []byte(ui.F1),
	glfw3.KeyF2:  []byte(ui.F2),
	glfw3.KeyF3:  []byte(ui.F3),
	glfw3.KeyF4:  []byte(ui.F4),
	glfw3.KeyF5:  []byte(ui.F5),
	glfw3.KeyF6:  []byte(ui.F6),
	glfw3.KeyF7:  []byte(ui.F7),
	glfw3.KeyF8:  []byte(ui.F8),
	glfw3.KeyF9:  []byte(ui.F9),
	glfw3.KeyF10: []byte(ui.F10),
	glfw3.KeyF11: []byte(ui.F11),
	glfw3.KeyF12: []byte(ui.F12),
	glfw3.KeyF13: []byte(ui.F13),
	glfw3.KeyF14: []byte(ui.F14),
	glfw3.KeyF15: []byte(ui.F15),
	glfw3.KeyF16: []byte(ui.F16),
	glfw3.KeyF17: []byte(ui.F17),
	glfw3.KeyF18: []byte(ui.F18),
	glfw3.KeyF19: []byte(ui.F19),
	glfw3.KeyF20: []byte(ui.F20),
	glfw3.KeyF21: []byte(ui.F21),
	glfw3.KeyF22: []byte(ui.F22),
	glfw3.KeyF23: []byte(ui.F23),
	glfw3.KeyF24: []byte(ui.F24),
	glfw3.KeyF25: []byte(ui.F25),
}

func translateKey(s []byte, k glfw3.Key) []byte {
	b, ok := keymap[k]
	if ok {
		return append(s, b...)
	}
	return s
}

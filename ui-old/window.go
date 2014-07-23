package ui

import (
	"errors"
	"fmt"
	"image"
	"reflect"

	"j4k.co/exp/ui/graphics"
	"j4k.co/exp/ui/scene"
)

type void struct{}

// Window is the root of the visual hierarchy in the user interface.
type Window struct {
	world      *scene.World
	viewport   *graphics.Viewport
	logics     map[reflect.Type]void
	prevm      Mouse
	controls   []*Control
	body       *Control
	mouseFocus *Control
}

// Open a new window using the specified environment.
func Open(env Environment) (*Window, error) {
	viewport := &graphics.Viewport{}
	world, err := scene.New(&graphics.System{
		Context:  env.Graphics(),
		Viewport: viewport,
	}, &controlsys{})
	if err != nil {
		return nil, err
	}
	w := &Window{
		world:    world,
		logics:   map[reflect.Type]void{},
		viewport: viewport,
	}
	w.init()
	go w.run(env)
	return w, nil
}

// run executes the window using env.
func (w *Window) run(env Environment) {
	keybd := env.Keyboard()
	mouse := env.Mouse()
	view := env.View()
	for {
		select {
		case k := <-keybd:
			//w.scene.keyboard(k)
			w.keyboard(k)
		case m := <-mouse:
			//w.scene.mouse(m)
			w.mouse(m)
			w.update()
		case v := <-view:
			//w.scene.view(v)
			println("wat")
			w.view(v)
			w.update()
		}
	}
}

func (w *Window) init() {
	w.Register(&body{})
	w.body = w.Insert(&Control{}, &body{})
	w.body.SetRect(image.Rect(0, 0, 500, 500))
	w.viewport.Width = 500
	w.viewport.Height = 500
}

func (w *Window) update() {
	dirty := make([]int, 0, len(w.controls)/2)
	for i := range w.controls {
		if w.controls[i].dirty {
			w.controls[i].dirty = false
			dirty = append(dirty, i)
		}
	}
	if len(dirty) > 0 {
		tx := w.world.Edit()
		defer w.world.Step()
		defer tx.Commit()
		for _, idx := range dirty {
			ctl := w.controls[idx]
			tx.Set(ctl.ref, control{
				Logic:  "main.dragbox",
				Rect:   ctl.rect,
				Parent: ctl.parent.ref,
				Kids:   controlRefs(ctl.kids),
			})
		}
	}
}

func controlRefs(ctls []*Control) []scene.Ref {
	refs := make([]scene.Ref, len(ctls))
	for i := range ctls {
		refs[i] = ctls[i].ref
	}
	return refs
}

func (w *Window) mouse(m Mouse) {
	if w.mouseFocus != nil && w.mouseFocus.focusLocked {
		logic, ok := w.mouseFocus.logic.(MouseListener)
		if ok {
			logic.Mouse(w.mouseFocus, m, w.prevm)
		}
	} else {
		ctl := w.hitTest(w.Body(), m.Point)
		if ctl != nil {
			if w.mouseFocus != ctl {
				w.switchMouse(ctl, w.mouseFocus)
			}
			logic, ok := ctl.logic.(MouseListener)
			if ok {
				logic.Mouse(ctl, m, w.prevm)
			}
		}
	}
	w.prevm = m
}

func (w *Window) switchMouse(ctl, prev *Control) {
	if prev != nil {
		logic, ok := prev.logic.(MouseFocusListener)
		if ok {
			logic.MouseFocus(ctl, false)
		}
	}
	if ctl != nil {
		logic, ok := ctl.logic.(MouseFocusListener)
		if ok {
			logic.MouseFocus(ctl, true)
		}
	}
	w.mouseFocus = ctl
}

func (w *Window) hitTest(ctl *Control, pt image.Point) *Control {
	for _, kid := range ctl.kids {
		tester, ok := kid.logic.(HitTester)
		if !ok {
			tester = defaultHitTester
		}
		if tester.HitTest(kid, pt) {
			return w.hitTest(kid, pt)
		}
	}
	return ctl
}

// keyboard sends keyboard input to the keyboard focused element.
func (w *Window) keyboard(k Keyboard) {
	//if w.keyfocus != nil {
	//w.keyfocus.Keyboard(k)
	//}
}

func (w *Window) view(v View) {
	fmt.Println(v)
	w.Body().SetRect(image.Rect(0, 0, v.Width, v.Height))
	//v.Ready()
	// TODO: server should just wait for swapbuffers
	//w.body.SetBounds(image.Rect(0, 0, v.Width, v.Height))
	//w.paint()
}

func (w *Window) Register(logics ...interface{}) {
	for _, l := range logics {
		w.register(l)
	}
}

func (w *Window) register(logic interface{}) {
	typ := reflect.TypeOf(logic)
	if typ.Kind() == reflect.Ptr {
		//typ = typ.Elem()
	}
	_, ok := w.logics[typ]
	if !ok {
		w.logics[typ] = void{}
	}
}

func typekey(typ reflect.Type) string {
	return fmt.Sprintf("%s.%s",
		typ.PkgPath(),
		typ.Name())
}

func (w *Window) createBody() {

}

// Body returns the root visual of the window.
func (w *Window) Body() *Control {
	return w.body
}

// SetParent sets visual's parent to parent. Panics if parent is nil or
// not a valid visual for the window. Orphans are not allowed.
func (w *Window) SetParent(ctl *Control, parent *Control) {
	if parent == nil {
		panic(errors.New("ui: nil parent"))
	}
	if ctl == w.Body() {
		panic(errors.New("ui: cannot set body parent"))
	}
	sibs := ctl.parent.kids
	for i := range sibs {
		if ctl == sibs[i] {
			copy(sibs[i:], sibs[i+1:len(sibs)-1])
			ctl.parent.kids = sibs[:len(sibs)-1]
			break
		}
	}
	parent.kids = append(parent.kids, ctl)
	ctl.parent = parent
	// update scene graph
}

// Insert a new visual into the scene using the specified logic and
// parent visual. Panics if parent is nil, or logic type is not
// registered.
func (w *Window) Insert(parent *Control, logic interface{}) *Control {
	if parent == nil && w.body != nil {
		panic(errors.New("ui: nil parent"))
	}
	typ := reflect.TypeOf(logic)
	if typ.Kind() == reflect.Ptr {
		//typ = typ.Elem()
	}
	_, ok := w.logics[typ]
	if !ok {
		panic(fmt.Errorf("ui: unknown logic type: %v", typ))
	}
	tx := w.world.Edit()
	defer tx.Commit()
	rect := image.Rect(0, 0, 100, 100)
	ctl := &Control{
		parent: parent,
		logic:  logic,
		rect:   rect,
		dirty:  true,
	}
	var pref scene.Ref
	if parent != nil {
		pref = parent.ref
		parent.kids = append(parent.kids, ctl)
		parent.dirty = true
	}
	ref := tx.Create(control{
		Parent: pref,
		Logic:  typekey(typ),
		Rect:   rect,
	})
	ctl.ref = ref
	w.controls = append(w.controls, ctl)
	return ctl
}

func (w *Window) UnmarshalJSON(data []byte) error {
	return nil
}

func (w *Window) MarshalJSON() ([]byte, error) {
	return nil, nil
}

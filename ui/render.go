// +build ignore

package ui

import (
	"github.com/go-gl/gl"
)

type Shader struct {
}

type Material struct {
	shader Shader
	path   string
}

type renderItem struct {
	Transform Transform
	Material  MaterialRef
	Geometry  string
}

type renderData struct {
	nodes map[NodeKey]int
	items []renderItem
}

type RenderSystem struct {
	frames chan func(Graphics)
}

func (r *RenderSystem) NodeEnter(node NodeKey, arch Archetype) {
}

func (r *RenderSystem) NodeExit(node NodeKey, arch Archetype) {
}

func (r *RenderSystem) Run(scene *Scene, sys SystemNode) {
	sys.Borrow(&layout)
	sys.Reads(&s.items)
	sys.Nodes(s.dirtyNodes)
	r.frames = make(chan func(Graphics), 1)
	go r.glthread(scene.graphics)
	for sys.NextFrame() {
		r.render()
	}
}

func (r *RenderSystem) render() {
	// TODO: acquire render items (letting subscribers access first)
	// TODO: defer release
	r.frames <- func(graphics Graphics) {
		r.draw()
		graphics.SwapBuffers()
	}
}

func (r *RenderSystem) draw() {

}

func (r *RenderSystem) glthread(graphics Graphics) {
	graphics.Lock()
	defer graphics.Unlock()
	for cmd := range r.cmds {
		cmd(graphics)
	}
}

func (r *RenderSystem) CreateItem(node NodeKey) RenderItemKey {
}

func (r *RenderSystem) DeleteItem(item RenderItemKey) {
}

func (r *RenderItem) SetItemTransform(item RenderItemKey, t Transform) {
}

func (r *RenderItem) SetItemMaterial(item RenderItemKey, m Material) {
}

func (r *RenderItem) SetItemGeometry(item RenderItemKey, g Geometry) {
}

////////////////

type renderable struct {
	// MVP transform
	// model transform
	// material (shader, render state, and texture resources)
	// geometry
}

func (r *renderable) render() {
}

type renderables []renderable

func (w *Window) renderInit() {
	//bounds := w.body.Bounds()
	gl.Disable(gl.MULTISAMPLE)
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	//gl.Viewport(0, 0, bounds.Max.X, bounds.Max.Y)
	gl.ClearDepth(1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)
}

func (w *Window) renderUpload() {
	/*
		for _, node := range w.nodes {
			node.gfxupload()
		}
	*/
	gl.Flush()
}

func (w *Window) prepareRenderList() {
	w.rlistmu.Lock()
	go func() {
		defer w.rlistmu.Unlock()
		rlist := w.rlist[:0]
		/*
			for _, node := range w.nodes {
				rlist = append(rlist, node.renderable())
			}
		*/
		//sort.Sort(renderables(rlist))
		w.rlist = rlist
	}()
}

func (w *Window) renderList() {
	w.rlistmu.Lock()
	defer w.rlistmu.Unlock()
	for _, r := range w.rlist {
		r.render()
	}
}

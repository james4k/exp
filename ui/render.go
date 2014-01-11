package ui

import (
	"github.com/go-gl/gl"
)

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
	bounds := w.body.Bounds()
	gl.Disable(gl.MULTISAMPLE)
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	gl.Viewport(0, 0, bounds.Max.X, bounds.Max.Y)
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
				rlist = append(rlist, node.renderables()...)
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

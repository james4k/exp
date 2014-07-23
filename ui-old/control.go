package ui

import (
	"image"
	"reflect"

	"j4k.co/exp/ui/graphics"
	"j4k.co/exp/ui/scene"
	"j4k.co/gfx"
	"j4k.co/gfx/geometry"
)

type Control struct {
	parent      *Control
	kids        []*Control
	logic       interface{}
	ref         scene.Ref
	rect        image.Rectangle
	dirty       bool
	focusLocked bool
}

func (c *Control) Rect() image.Rectangle {
	return c.rect
}

func (c *Control) SetRect(rect image.Rectangle) {
	c.rect = rect
	c.dirty = true
}

func (c *Control) SetFocusLock(locked bool) {
	c.focusLocked = locked
}

// control represents a ui control in the scene graph.
type control struct {
	Logic  string
	Rect   image.Rectangle
	Parent scene.Ref
	Kids   []scene.Ref
}

type controlColumn struct {
	ctls   map[scene.Ref]control
	meshes map[scene.Ref]scene.Ref
	dirty  map[scene.Ref]struct{}
}

func (c *controlColumn) Type() interface{} {
	return reflect.TypeOf((*control)(nil)).Elem()
}

func (c *controlColumn) Add(node scene.Ref, data interface{}) {
	ctl := data.(control)
	ctl.Rect = image.Rect(0, 0, 100, 100)
	c.ctls[node] = ctl
	c.mark(node)
}

func (c *controlColumn) Del(node scene.Ref) {
	delete(c.ctls, node)
}

func (c *controlColumn) Get(node scene.Ref, data interface{}) {
}

func (c *controlColumn) Set(node scene.Ref, data interface{}) {
	c.ctls[node] = data.(control)
	c.mark(node)
}

func (c *controlColumn) mark(node scene.Ref) {
	c.dirty[node] = struct{}{}
}

type controlsys struct {
}

func (c *controlsys) Run(ctx *scene.Context) {
	var (
		ctls   controlColumn
		meshes graphics.MeshColumn
	)
	ctx.Stage(scene.DefaultStage - 1)
	ctx.Declare(&controlColumn{
		ctls:   map[scene.Ref]control{},
		meshes: map[scene.Ref]scene.Ref{},
		dirty:  map[scene.Ref]struct{}{},
	})
	ctx.Require(&ctls, &meshes)
	for ctx.Step() {
		dirty := make([]scene.Ref, len(ctls.dirty))
		for ref := range ctls.dirty {
			ctl := ctls.ctls[ref]
			meshref := ctls.meshes[ref]
			if meshref == 0 {
				meshref = ctx.Alloc(&meshes)
				gfx.DefaultVertexAttributes[gfx.VertexPosition] = "Position"
				gfx.DefaultVertexAttributes[gfx.VertexColor] = "Color"
				gfx.DefaultVertexAttributes[gfx.VertexTexcoord] = "UV"
				// diffuse = white texture
				// color = white, but by default i guess?
				meshes.Add(meshref, graphics.Mesh{
					G: graphics.Geometry{geometry.NewBuilder(gfx.DefaultVertexAttributes.Format())},
					M: graphics.M().Diffuse(0),
				})
				ctls.meshes[ref] = meshref
			}
			g := meshes.ModifyGeom(meshref)
			buildRect(g, ctl.Rect)
			dirty = append(dirty, ref)
		}
		for _, ref := range dirty {
			delete(ctls.dirty, ref)
		}
	}
}

func buildRect(g graphics.Geometry, rect image.Rectangle) {
	g.Clear()
	//g.SetCapacity(4, 6)
	g.P(float32(rect.Min.X), float32(rect.Min.Y), 0)
	g.UV(0, 0).C(255, 255, 255, 180)
	g.P(float32(rect.Max.X), float32(rect.Min.Y), 0)
	g.UV(1, 0)
	g.P(float32(rect.Max.X), float32(rect.Max.Y), 0)
	g.UV(1, 1)
	g.P(float32(rect.Min.X), float32(rect.Max.Y), 0)
	g.UV(0, 1)
	g.Indices(0, 1, 2, 2, 0, 3)
}

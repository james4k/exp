package graphics

import (
	"reflect"

	"j4k.co/exp/ui/scene"
	"j4k.co/gfx"
	"j4k.co/gfx/geometry"
)

type TextureRef scene.Ref

// Geometry is a geometry builder that can be serialized by the scene.
type Geometry struct {
	*geometry.Builder
}

// A Mesh is a Geometry-Material pair.
type Mesh struct {
	G Geometry
	M Material
}

type MeshColumn struct {
	m       map[scene.Ref]Mesh
	geoms   map[scene.Ref]gfx.Geometry
	layouts map[scene.Ref]gfx.GeometryLayout
	dirty   map[scene.Ref]struct{}
	typ     reflect.Type
}

func meshCol() *MeshColumn {
	return &MeshColumn{
		m:       map[scene.Ref]Mesh{},
		geoms:   map[scene.Ref]gfx.Geometry{},
		layouts: map[scene.Ref]gfx.GeometryLayout{},
		dirty:   map[scene.Ref]struct{}{},
	}
}

func (c *MeshColumn) Type() interface{} {
	return reflect.TypeOf((*Mesh)(nil)).Elem()
}

func (c *MeshColumn) Add(node scene.Ref, data interface{}) {
	c.m[node] = data.(Mesh)
}

func (c *MeshColumn) Del(node scene.Ref) {
	delete(c.m, node)
}

func (c *MeshColumn) Get(node scene.Ref, data interface{}) {
	_, ok := c.m[node]
	if ok {
		// TODO: set via reflection.. data is a ptr
	}
}

func (c *MeshColumn) Set(node scene.Ref, data interface{}) {
	c.m[node] = data.(Mesh)
}

// ModifyGeom returns the geometry buffer for node, and marks it as dirty.
func (c *MeshColumn) ModifyGeom(node scene.Ref) Geometry {
	mesh, ok := c.m[node]
	if !ok {
		return Geometry{nil}
	}
	c.dirty[node] = struct{}{}
	return mesh.G
}

func (c *MeshColumn) update(cmds glchan, shader *gfx.Shader) {
	var newrefs []scene.Ref
	dirty := make([]scene.Ref, 0, len(c.dirty))
	geoms := make([]gfx.Geometry, 0, len(c.dirty))
	data := make([]Geometry, 0, len(c.dirty))
	newdata := make([]Geometry, 0, len(c.dirty))
	for ref := range c.dirty {
		dirty = append(dirty, ref)
		geom, ok := c.geoms[ref]
		if ok {
			geoms = append(geoms, geom)
			data = append(data, c.m[ref].G)
		} else {
			newrefs = append(newrefs, ref)
			newdata = append(newdata, c.m[ref].G)
		}
	}
	for _, ref := range dirty {
		delete(c.dirty, ref)
	}
	newgeoms := make([]gfx.Geometry, len(newrefs))
	allocGeom(cmds, newgeoms)
	for i, ref := range newrefs {
		c.geoms[ref] = newgeoms[i]
	}
	copyGeom(cmds, geoms, data)
	copyGeom(cmds, newgeoms, newdata)
	newlayouts := make([]gfx.GeometryLayout, len(newgeoms))
	bindGeom(cmds, newlayouts, newgeoms, shader)
	for i := range newlayouts {
		c.layouts[newrefs[i]] = newlayouts[i]
	}
}

func allocGeom(cmds glchan, dst []gfx.Geometry) {
	if len(dst) == 0 {
		return
	}
	ch := make(chan struct{})
	cmds <- func() {
		gfx.AllocGeometry(dst, gfx.DynamicDraw)
		close(ch)
	}
	<-ch
}

func copyGeom(cmds glchan, dst []gfx.Geometry, data []Geometry) {
	if len(data) == 0 {
		return
	}
	cmds <- func() {
		for i := range dst {
			dst[i].CopyFrom(data[i])
		}
	}
}

func bindGeom(cmds glchan, dst []gfx.GeometryLayout, geom []gfx.Geometry, shader *gfx.Shader) {
	ch := make(chan struct{})
	cmds <- func() {
		shader.Use()
		for i := range geom {
			dst[i].Layout(shader, &geom[i])
		}
		close(ch)
	}
	<-ch
}

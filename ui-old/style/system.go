package style

import (
	"reflect"

	"j4k.co/exp/ui/scene"
	"j4k.co/gfx"
)

type Prop uint16

const (
	BackgroundColor Prop = iota
)

type Prop struct {
	Name  string
	Value string
}

type Properties struct {
	Props []Prop
}

type propertyColumn struct {
	m       map[scene.Ref]Mesh
	geoms   map[scene.Ref]gfx.Geometry
	layouts map[scene.Ref]gfx.GeometryLayout
	dirty   map[scene.Ref]struct{}
	typ     reflect.Type
}

func propCol() *propertyColumn {
	return &propertyColumn{
		m:       map[scene.Ref]Mesh{},
		geoms:   map[scene.Ref]gfx.Geometry{},
		layouts: map[scene.Ref]gfx.GeometryLayout{},
		dirty:   map[scene.Ref]struct{}{},
	}
}

func (c *propertyColumn) Type() interface{} {
	return reflect.TypeOf((*Mesh)(nil)).Elem()
}

func (c *propertyColumn) Add(node scene.Ref, data interface{}) {
	c.m[node] = data.(Mesh)
}

func (c *propertyColumn) Del(node scene.Ref) {
	delete(c.m, node)
}

func (c *propertyColumn) Get(node scene.Ref, data interface{}) {
	_, ok := c.m[node]
	if ok {
		// TODO: set via reflection.. data is a ptr
	}
}

func (c *propertyColumn) Set(node scene.Ref, data interface{}) {
	c.m[node] = data.(Mesh)
}

type System struct {
}

func (s *System) Run(ctx *scene.Context) {
	ctx.Stage(scene.StyleStage)
	ctx.Declare(meshCol())
	var (
		meshes MeshColumn
	)
	ctx.Require(&meshes)

	for ctx.Step() {
	}
}

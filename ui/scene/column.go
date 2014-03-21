package scene

import "reflect"

// world should own a column by default, until a system wants it.

// Column is a collection of data of the same type.
type Column interface {
	Type() interface{}
	Add(Ref, interface{})
	Del(Ref)
	Get(Ref, interface{})
	Put(Ref, interface{})
}

type column struct {
	m   map[Ref]interface{}
	typ reflect.Type
}

func DumbColumn(zeroval interface{}) Column {
	return column{
		m:   map[Ref]interface{}{},
		typ: reflect.TypeOf(zeroval),
	}
}

func (c column) Type() interface{} {
	return c.typ
}

func (c column) Add(node Ref, data interface{}) {
	c.m[node] = data
}

func (c column) Del(node Ref) {
	delete(c.m, node)
}

func (c column) Get(node Ref, data interface{}) {
	_, ok := c.m[node]
	if ok {
		// TODO: set via reflection.. data is a ptr
	}
}

func (c column) Put(node Ref, data interface{}) {
	// TODO: check if data is indirect and dereference
	c.m[node] = data
}

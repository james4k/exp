package graphics

// So, we want to clean this up to where.. material fields go through an
// interface? Look at vertex attributes? Well, let's look at what
// shaders need.

/*
type Material interface {
	Bool(FieldKey) bool
	Color(FieldKey) image.Color
}
*/

type Material struct {
	fields []field
}

//var Default = Material{}

func M() Material {
	return Material{}
}

func (m Material) Diffuse(ref TextureRef) Material {
	return m.put(textureField{
		fieldKey: diffuseKey,
		Ref:      ref,
	})
}

func (m Material) put(field field) Material {
	key := field.key()
	for i, f := range m.fields {
		if f.key() == key {
			m.fields[i] = field
			return m
		}
	}
	m.fields = append(m.fields, field)
	return m
}

func (m Material) field(key fieldKey) field {
	for _, f := range m.fields {
		if f.key() == key {
			return f
		}
	}
	return nil
}

type fieldKey uint32

func (k fieldKey) key() fieldKey {
	return k
}

const (
	diffuseKey fieldKey = 1 + iota
)

type field interface {
	key() fieldKey
}

type textureField struct {
	fieldKey
	Ref TextureRef
}

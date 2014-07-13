package graphics

import (
	"image"
	"image/color"
	"os"
	"sync"

	"github.com/Jragonmiris/mathgl"
	"github.com/go-gl/gl"
	"j4k.co/exp/ui/scene"
	"j4k.co/gfx"
	"j4k.co/gfx/geometry"
)

type glchan chan<- func()

type Viewport struct {
	Width, Height int
}

const (
	opaqueBucket = iota
	transparentBucket
	numBuckets
)

type System struct {
	Context
	*Viewport
	shader  gfx.Shader
	buckets []renderList

	stepmu sync.Mutex
	cmds   chan func()
	meshes MeshColumn
}

func (s *System) init(ctx *scene.Context) {
	s.cmds = make(chan func(), 4)
	go glthread(s.Context, s.cmds)
	s.cmds <- s.glsetup
	s.stepmu.Lock()
	ctx.Stage(scene.OutputStage)
	ctx.Declare(meshCol())
	ctx.Require(&s.meshes)
	s.buckets = make([]renderList, numBuckets)
}

func (s *System) finish() {
	close(s.cmds)
}

func (s *System) Run(ctx *scene.Context) {
	s.init(ctx)
	defer s.finish()
	// TODO: take active camera node
	w := float32(s.Width)
	h := float32(s.Height)
	projM := mathgl.Ortho(0, w, h, 0, -30, 30)
	viewM := mathgl.Translate3D(0, 0, 0)
	worldM := mathgl.Translate3D(0, 0, 0)
	var rect struct {
		WVP     mathgl.Mat4f   `uniform:"WorldViewProjectionM"`
		Diffuse *gfx.Sampler2D `uniform:"Diffuse"`
		geomobj gfx.Geometry
		geom    gfx.GeometryLayout
	}
	rect.WVP = projM.Mul4(viewM).Mul4(worldM)
	s.cmds <- func() {
		quadbuffer := geometry.NewBuilder(s.shader.VertexFormat())
		quadbuffer.Clear()
		quadbuffer.P(0, 0, 0).UV(-1, -1).Cf(1, 1, 1, 1)
		quadbuffer.P(w/3, 0, 0).UV(-1, -1)
		quadbuffer.P(w/3, h/3, 0).UV(-1, -1)
		quadbuffer.P(0, h/3, 0).UV(-1, -1)
		quadbuffer.Indices(0, 1, 2, 2, 0, 3)
		rect.geomobj.Alloc(gfx.StaticDraw)
		err := rect.geomobj.CopyFrom(quadbuffer)
		if err != nil {
			panic(err)
		}
		s.shader.Use()
		rect.geom.Layout(&s.shader, &rect.geomobj)
		/*
			err = rect.geomobj.CopyFrom(quadbuffer)
			if err != nil {
				panic(err)
			}
		*/
		whiteImg := image.NewNRGBA(image.Rect(0, 0, 1, 1))
		whiteImg.Set(0, 0, color.White)
		white, err := gfx.Image(whiteImg)
		if err != nil {
			panic(err)
		}
		rect.Diffuse = white
	}

	for ctx.Step() {
		s.cmds <- func() {
			gl.Viewport(0, 0, s.Width, s.Height)
			gl.ClearColor(0, 0, 0, 1.0)
			gl.ClearDepth(1)
			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		}
		s.meshes.update(s.cmds, &s.shader)
		s.draw()
		s.sync()
	}
}

func (s *System) sync() {
	s.cmds <- s.stepmu.Unlock
	s.cmds <- s.SwapBuffers
	s.stepmu.Lock()
}

func (s *System) draw() {
	s.cmds <- func() {
		gl.Enable(gl.DEPTH_TEST)
		gl.Enable(gl.BLEND)
		gl.DepthFunc(gl.LEQUAL)
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
		//s.shader.Use()
		//s.shader.AssignUniforms(uniforms)
	}
	for _, _ = range s.meshes.layouts {
		s.submit(renderState{})
	}
	for i := range s.buckets {
		list := &s.buckets[i]
		list.sort()
		s.cmds <- list.execute
	}
	return
	/*
		s.cmds <- func() {
			for _, geom := range geoms {
				// gl state
				// texture units
				// uniforms
				//shader.AssignUniforms()
				shader.SetGeometry(&geom)
				shader.Draw()
			}
		}
	*/
}

func (s *System) submit(state renderState) {
	key := state.key()
	for i := range s.buckets {
		//if s.buckets[i].match(key) {
		//}
		s.buckets[i].insert(key, state)
		break
	}
}

func (s *System) glsetup() {
	gfx.DefaultVertexAttributes[gfx.VertexPosition] = "Position"
	gfx.DefaultVertexAttributes[gfx.VertexColor] = "Color"
	gfx.DefaultVertexAttributes[gfx.VertexTexcoord] = "UV"
	//gfx.DefaultVertexAttributes[gfx.VertexNormal] = "Normal"
	s.shader.Build(os.Stderr, gfx.DefaultVertexAttributes, fs, vs)
}

func glthread(glctx Context, cmds <-chan func()) {
	glctx.Lock()
	defer glctx.Unlock()
	gl.Init()
	for cmd := range cmds {
		cmd()
	}
}

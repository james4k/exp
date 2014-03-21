package scene

import "reflect"

// Stage specifies when in the world step a system will process its
// dependencies. Subsystems in the same stage may run concurrently if
// they do not share dependencies, and order is undefined.
type Stage int

const (
	DefaultStage Stage = 0   // most work is done in the default stage
	OutputStage        = 100 // scene is rendered to screen
)

// Subsystem represents an independent process in the world, and runs in
// its own goroutine. Its execution is synchronized in lockstep with all
// other systems in the world for frame consistency, while step order is
// determined by its stage and its data dependencies.
type Subsystem interface {
	Run(*Context)
}

// Context is a subsystem's context of a world.
type Context struct {
	stage   *stageClient
	stageId Stage
	world   *World
	req     columnReq
}

// Stage sets which stage this system will run in. Panics if called
// after Step has been called.
func (c *Context) Stage(s Stage) {
	if c.stage != nil {
		panic("scene: context already initialized")
	}
	c.stageId = s
}

// Declare takes an initialized column to be added to the world. Order
// specifies when the system will process the column in relation to
// other systems. The system should not hold references to columns, and
// should be acquired as dependencies in the Step method. Panics if a
// column of the same type already exists or if Step has been called.
func (c *Context) Declare(columns ...Column) {
	if c.stage != nil {
		panic("scene: context already initialized")
	}
	c.world.declare(columns)
}

// Require takes in pointers to column values that the subsystem depends
// on for processing. These values will be set after calling Step and
// are only valid until Step is reached again.
func (c *Context) Require(dependencies ...Column) {
	for _, col := range dependencies {
		val := reflect.ValueOf(col)
		c.req.dest = append(c.req.dest, val.Elem())
	}
}

// Step returns true when the world is ready to take another step. When
// the world ends, Step returns false.
func (c *Context) Step() bool {
	if c.stage == nil {
		c.world.ctxInit(c)
	}
	c.release()
	if !c.stage.wait() {
		return false
	}
	c.acquire()
	return true
}

// release dependencies for other subsystems to use
func (c *Context) release() {
	c.req.copyBack()
	c.world.release(&c.req)
}

// acquire exclusive access to dependencies defined by calls to Require
func (c *Context) acquire() {
	c.world.acquire(&c.req)
	c.req.copy()
}

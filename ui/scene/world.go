package scene

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"sync"
)

type void struct{}

// columnReq defines a set of columns to be acquired exclusively at
// once to avoid AB-BA deadlocks as well as better potential
// parallelism.
type columnReq struct {
	slots []*columnSlot
	dest  []reflect.Value
}

func (c *columnReq) copy() {
	for i := range c.slots {
		c.dest[i].Set(c.slots[i].val)
	}
}

func (c *columnReq) copyBack() {
	for i := range c.slots {
		c.slots[i].val.Set(c.dest[i])
	}
}

// columnSlot is where column requests are fulfilled from
type columnSlot struct {
	cond   *sync.Cond
	val    reflect.Value
	locked bool
}

// World represents a scene.
type World struct {
	initc     chan *Context
	initdone  chan void
	mu        sync.Mutex
	nextref   Ref
	nodes     map[Ref]Column
	colByType map[reflect.Type]Column
	columns   []Column
	slotmu    sync.Mutex
	colslots  []columnSlot
	systems   []Subsystem
	contexts  []Context
	stages    []stage
}

func New(systems ...Subsystem) (*World, error) {
	w := &World{
		nodes:     map[Ref]Column{},
		colByType: map[reflect.Type]Column{},
	}
	err := w.init(systems...)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (w *World) init(systems ...Subsystem) error {
	w.systems = make([]Subsystem, len(systems))
	w.contexts = make([]Context, len(systems))
	copy(w.systems, systems)
	w.initc = make(chan *Context)
	w.initdone = make(chan void)
	for i, sys := range systems {
		ctx := &w.contexts[i]
		ctx.world = w
		go w.run(sys, ctx)
	}
	return w.initContexts()
}

func recoverError(err *error) {
	x := recover()
	switch x := x.(type) {
	case error:
		*err = x
	case string:
		*err = errors.New(x)
	}
}

func (w *World) initContexts() (rerr error) {
	defer recoverError(&rerr)
	stageIds := map[Stage]void{}
	for _ = range w.contexts {
		ctx := <-w.initc
		stageIds[ctx.stageId] = void{}
	}
	w.stages = make([]stage, len(stageIds))
	idx := 0
	for sid := range stageIds {
		w.stages[idx].init(sid)
		idx++
	}
	sort.Sort(stagesById(w.stages))
	for i := range w.contexts {
		ctx := &w.contexts[i]
		idx = w.stageIndex(ctx.stageId)
		if idx < 0 {
			panic("scene: subsystem init failed")
		}
		w.stages[idx].add()
		ctx.stage = &w.stages[idx].stageClient
	}
	for i := range w.contexts {
		ctx := &w.contexts[i]
		ctx.req.slots = make([]*columnSlot, len(ctx.req.dest))
		for i := range ctx.req.dest {
			ctx.req.slots[i] = w.findSlot(ctx.req.dest[i])
		}
		ctx.req.copy()
	}
	close(w.initdone)
	w.initc = nil
	// TODO: we would nil w.initdone but that is not easy to do without
	// being racy. ideally this would all be on the stack
	return rerr
}

func (w *World) findSlot(val reflect.Value) *columnSlot {
	typ := val.Type()
	for i := range w.colslots {
		if typ == w.colslots[i].val.Type() {
			return &w.colslots[i]
		}
	}
	panic(fmt.Errorf("scene: slot not found for %v", typ))
}

func (w *World) ctxInit(ctx *Context) {
	w.initc <- ctx
	<-w.initdone
}

func (w *World) run(sys Subsystem, ctx *Context) {
	//defer
	sys.Run(ctx)
}

func (w *World) stageIndex(id Stage) int {
	for i := range w.stages {
		if w.stages[i].id == id {
			return i
		}
	}
	return -1
}

func (w *World) Step() {
	w.stop()
	defer w.start()
	for i := range w.stages {
		w.stages[i].cycle()
	}
}

func (w *World) Close() {
	w.stop()
	defer w.start()
	for i := range w.stages {
		w.stages[i].kill()
	}
}

func (w *World) stop() {
	w.mu.Lock()
}

func (w *World) start() {
	w.mu.Unlock()
}

func (w *World) declare(columns []Column) {
	w.stop()
	defer w.start()
	for _, col := range columns {
		typ := col.Type().(reflect.Type)
		if typ == nil {
			panic(fmt.Errorf("scene: column %T returns invalid type", col))
		}
		_, ok := w.colByType[typ]
		if ok {
			panic(fmt.Errorf("scene: column for %v already declared", typ))
		}
		w.colByType[typ] = col
		val := reflect.ValueOf(col)
		if val.Kind() != reflect.Ptr {
			// TODO: perhaps support slices and maps
			newval := reflect.New(val.Type()).Elem()
			newval.Set(val)
			val = newval
		} else {
			val = val.Elem()
		}
		w.columns = append(w.columns, col)
		w.colslots = append(w.colslots, columnSlot{
			cond: sync.NewCond(&w.slotmu),
			val:  val,
		})
	}
}

func (w *World) acquire(req *columnReq) {
	if len(req.slots) == 0 {
		return
	}
	w.slotmu.Lock()
	defer w.slotmu.Unlock()
loop:
	for {
		// only lock each slot if we can lock all of them
		for _, slot := range req.slots {
			if slot.locked {
				// w.slotmu unlocked by cond wait
				slot.cond.Wait()
				// w.slotmu locked again on awakening
				continue loop
			}
		}
		for _, slot := range req.slots {
			slot.locked = true
		}
		return
	}
}

func (w *World) release(req *columnReq) {
	if len(req.slots) == 0 {
		return
	}
	w.slotmu.Lock()
	defer w.slotmu.Unlock()
	for _, slot := range req.slots {
		slot.locked = false
		slot.cond.Broadcast()
	}
}

func (w *World) Edit() Tx {
	w.stop()
	return Tx{w}
}

/*
func (w *World) Create(data interface{}) Ref {
	w.stop()
	defer w.start()
	return w.create(data)
}

func (w *World) Delete(node Ref) {
	w.stop()
	defer w.start()
	w.delete(node)
}
*/

func (w *World) create(data interface{}) Ref {
	typ := reflect.TypeOf(data)
	c := w.colByType[typ]
	if c == nil {
		for ty, col := range w.colByType {
			fmt.Printf("%T - %T\n", ty, col)
		}
		panic(fmt.Errorf("scene: no column for type %T", data))
	}
	ref := w.alloc(c)
	c.Add(ref, data)
	return ref
}

func (w *World) alloc(c Column) Ref {
	for {
		_, ok := w.nodes[w.nextref]
		if !ok && w.nextref != Nil {
			break
		}
		w.nextref++
	}
	ref := w.nextref
	w.nextref++
	w.nodes[ref] = c
	return ref
}

func (w *World) delete(node Ref) {
}

func (w *World) set(node Ref, data interface{}) {
	col := w.nodes[node]
	if col != nil {
		col.Set(node, data)
	}
}

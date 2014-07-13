package scene_test

import (
	"testing"

	"j4k.co/exp/ui/scene"
)

type sysFunc func(*scene.Context)

func (s sysFunc) Run(ctx *scene.Context) {
	s(ctx)
}

func TestCycleOrder(t *testing.T) {
	type clockColumn struct {
		scene.Column
		clock [3]int
	}
	sysA := sysFunc(func(ctx *scene.Context) {
		ctx.Declare(clockColumn{
			Column: scene.DumbColumn(0),
		})
		var c clockColumn
		ctx.Require(&c)
		for ctx.Step() {
			c.clock[0]++
			switch {
			case c.clock[0] != c.clock[1]+1,
				c.clock[0] != c.clock[2]+1:
				panic("system out of stage")
			}
		}
	})
	sysB := sysFunc(func(ctx *scene.Context) {
		var c clockColumn
		ctx.Stage(1)
		ctx.Require(&c)
		for ctx.Step() {
			c.clock[1]++
			switch {
			case c.clock[1] != c.clock[0],
				c.clock[1] != c.clock[2]+1:
				panic("system out of stage")
			}
		}
	})
	sysC := sysFunc(func(ctx *scene.Context) {
		var c clockColumn
		ctx.Stage(2)
		ctx.Require(&c)
		for ctx.Step() {
			c.clock[2]++
			switch {
			case c.clock[2] != c.clock[0],
				c.clock[2] != c.clock[1]:
				panic("system out of stage")
			}
		}
	})
	world, err := scene.New(sysA, sysB, sysC)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		world.Step()
	}
	world.Close()
}

func TestSubsystemConcurrency(t *testing.T) {
	const N = 10
	numc := make(chan int)
	sysA := sysFunc(func(ctx *scene.Context) {
		count := 0
		for ctx.Step() {
			count++
			numc <- count
		}
	})
	sysB := sysFunc(func(ctx *scene.Context) {
		count := 0
		for ctx.Step() {
			count = <-numc
		}
		if count != N {
			panic("mismatch")
		}
	})
	world, err := scene.New(sysA, sysB)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < N; i++ {
		world.Step()
	}
	world.Close()
}

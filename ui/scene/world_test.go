package scene_test

import (
	"fmt"
	"reflect"
	"testing"

	"j4k.co/exp/ui/scene"
)

type clockColumn struct {
	scene.Column
	clock [3]int
}

type sysA struct {
}

func (s *sysA) Run(ctx *scene.Context) {
	var clock clockColumn
	ctx.Declare(clockColumn{
		Column: scene.DumbColumn(reflect.TypeOf(1)),
		clock:  [3]int{10, 0, 0},
	})
	ctx.Require(&clock)
	for ctx.Step() {
		clock.clock[0]++
		fmt.Println(clock.clock[0])
		// do work on data
	}
}

type sysB struct {
}

func (s *sysB) Run(ctx *scene.Context) {
	for ctx.Step() {
	}
}

type sysC struct {
}

func (s *sysC) Run(ctx *scene.Context) {
	for ctx.Step() {
	}
}

func TestCycleOrder(t *testing.T) {
	world, err := scene.New(&sysA{}, &sysB{}, &sysC{})
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		world.Step()
	}
}
